package dynamo

import (
	"encoding/json"
	"time"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	//"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/contamobi/oauth2"
	"github.com/contamobi/oauth2/models"
)

func getTokenTableDefinition(
	tableName string, 
	txnName string, 
	accessToken string, 
	refreshToken string,
	basic string,
	) *dynamodb.CreateTableInput {
	params := &dynamodb.CreateTableInput{
		TableName: tableName,
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: txnName,
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: accessToken,
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: refreshToken,
				AttributeName: aws.String("S"),
			},
			{
				AttributeName: basic,
				AtrributeName: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: accessToken,
				KeyType: aws.String("HASH"),
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughtput{
			ReadCapacityUnits: aws.Int64(10),
			WriteCapacityUnits: aws.Int64(5),
		},
	}

	return params
}

func tableHasExits(tableName string, svc *dynamodb.DynamoDB) bool {
	describeInput := &dynamodb.DescribeTableInput{
		TableName: aws.String(tableName),
	}
	describeTableOutput, err := svc.DescribeTable(describeInput)

	if err != nil {
		fmt.Printf("%s", err)
		return true
	}
	return false
}

func NewTokenStore(config *Config) (tokenStore *TokenStore, err error){
	session, err := config.SESSION
	if err != nil {
		return err
	}
	svc := dynamodb.New(session)
	if tableHasExits(config.TABLE, svc) != false {
		tableDefinition := getTokenTableDefinition(config.TABLE, "oauth2_txn", "oauth2_access", "oauth2_refresh", "oauth2_basic")
		_, err := svc.CreateTable(tableDefinition)
		if err != nil {
			return err
		}
	}
	ts := &TokenStore{
		config: config,
		session: svc,
	}
	tokenStorage = ts
	return
}

type TokenStore struct {
	config  *Config
	session *dynamodb.DynamoDB
}

// Create create and store the new token information
func (tokenStorage *TokenStore) Create(info oauth2.TokenInfo) (err error) {
	jv, err := json.Marshal(info)
	if err != nil {
		return
	}

	if code := info.GetCode(); code != "" {
		ts.cHandler(ts.tcfg.BasicCName, func(c *mgo.Collection) {
			err = c.Insert(basicData{
				ID:        code,
				Data:      jv,
				ExpiredAt: info.GetCodeCreateAt().Add(info.GetCodeExpiresIn()),
			})
		})
		return
	}

	aexp := info.GetAccessCreateAt().Add(info.GetAccessExpiresIn())
	rexp := aexp
	if refresh := info.GetRefresh(); refresh != "" {
		rexp = info.GetRefreshCreateAt().Add(info.GetRefreshExpiresIn())
		if aexp.Second() > rexp.Second() {
			aexp = rexp
		}
	}
	id := bson.NewObjectId().Hex()
	ops := []txn.Op{{
		C:      ts.tcfg.BasicCName,
		Id:     id,
		Assert: txn.DocMissing,
		Insert: basicData{
			Data:      jv,
			ExpiredAt: rexp,
		},
	}, {
		C:      ts.tcfg.AccessCName,
		Id:     info.GetAccess(),
		Assert: txn.DocMissing,
		Insert: tokenData{
			BasicID:   id,
			ExpiredAt: aexp,
		},
	}}
	if refresh := info.GetRefresh(); refresh != "" {
		ops = append(ops, txn.Op{
			C:      ts.tcfg.RefreshCName,
			Id:     refresh,
			Assert: txn.DocMissing,
			Insert: tokenData{
				BasicID:   id,
				ExpiredAt: rexp,
			},
		})
	}
	ts.cHandler(ts.tcfg.TxnCName, func(c *mgo.Collection) {
		runner := txn.NewRunner(c)
		err = runner.Run(ops, "", nil)
	})
	return
}

// RemoveByCode use the authorization code to delete the token information
func (ts *TokenStore) RemoveByCode(code string) (err error) {
	ts.cHandler(ts.tcfg.BasicCName, func(c *mgo.Collection) {
		verr := c.RemoveId(code)
		if verr != nil {
			if verr == mgo.ErrNotFound {
				return
			}
			err = verr
		}
	})
	return
}

// RemoveByAccess use the access token to delete the token information
func (ts *TokenStore) RemoveByAccess(access string) (err error) {
	ts.cHandler(ts.tcfg.AccessCName, func(c *mgo.Collection) {
		verr := c.RemoveId(access)
		if verr != nil {
			if verr == mgo.ErrNotFound {
				return
			}
			err = verr
		}
	})
	return
}

// RemoveByRefresh use the refresh token to delete the token information
func (ts *TokenStore) RemoveByRefresh(refresh string) (err error) {
	ts.cHandler(ts.tcfg.RefreshCName, func(c *mgo.Collection) {
		verr := c.RemoveId(refresh)
		if verr != nil {
			if verr == mgo.ErrNotFound {
				return
			}
			err = verr
		}
	})
	return
}

func (ts *TokenStore) getData(basicID string) (ti oauth2.TokenInfo, err error) {
	ts.cHandler(ts.tcfg.BasicCName, func(c *mgo.Collection) {
		var bd basicData
		verr := c.FindId(basicID).One(&bd)
		if verr != nil {
			if verr == mgo.ErrNotFound {
				return
			}
			err = verr
			return
		}
		var tm models.Token
		err = json.Unmarshal(bd.Data, &tm)
		if err != nil {
			return
		}
		ti = &tm
	})
	return
}

func (ts *TokenStore) getBasicID(cname, token string) (basicID string, err error) {
	ts.cHandler(cname, func(c *mgo.Collection) {
		var td tokenData
		verr := c.FindId(token).One(&td)
		if verr != nil {
			if verr == mgo.ErrNotFound {
				return
			}
			err = verr
			return
		}
		basicID = td.BasicID
	})
	return
}

// GetByCode use the authorization code for token information data
func (ts *TokenStore) GetByCode(code string) (ti oauth2.TokenInfo, err error) {
	ti, err = ts.getData(code)
	return
}

// GetByAccess use the access token for token information data
func (ts *TokenStore) GetByAccess(access string) (ti oauth2.TokenInfo, err error) {
	basicID, err := ts.getBasicID(ts.tcfg.AccessCName, access)
	if err != nil && basicID == "" {
		return
	}
	ti, err = ts.getData(basicID)
	return
}

// GetByRefresh use the refresh token for token information data
func (ts *TokenStore) GetByRefresh(refresh string) (ti oauth2.TokenInfo, err error) {
	basicID, err := ts.getBasicID(ts.tcfg.RefreshCName, refresh)
	if err != nil && basicID == "" {
		return
	}
	ti, err = ts.getData(basicID)
	return
}

type basicData struct {
	ID        string    `bson:"_id"`
	Data      []byte    `bson:"Data"`
	ExpiredAt time.Time `bson:"ExpiredAt"`
}

type tokenData struct {
	ID        string    `bson:"_id"`
	BasicID   string    `bson:"BasicID"`
	ExpiredAt time.Time `bson:"ExpiredAt"`
}
