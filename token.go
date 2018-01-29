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

func NewTokenStore(config *Config) (tokenStore *TokenStore) {
	session := config.SESSION
	svc := dynamodb.New(session)
	return &TokenStore{
		config: config,
		session: svc,
	}
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
		_,err := CreateWithAuthorizationCode(info)
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

func CreateWithAuthorizationCode(tokenStorage *TokenStore, info oauth2.TokenInfo) (info oauth2.TokenInfo) (err error){
	
	params := tokenStorage.session.PutItemInput{
		Item: vv,
		TableName: aws.String("table name")
	}
	err = c.Insert(basicData{
		ID:        code,
		Data:      jv,
		ExpiredAt: info.GetCodeCreateAt().Add(info.GetCodeExpiresIn()),
	})
}

func CreateWithAccessToken(info oauth2.TokenInfo) {

}

func CreateWithRefreshToken(info oauth2.TokenInfo) {

}

// RemoveByCode use the authorization code to delete the token information
func (tokenStorage *TokenStore) RemoveByCode(code string) (err error) {
	
}

// RemoveByAccess use the access token to delete the token information
func (tokenStorage *TokenStore) RemoveByAccess(access string) (err error) {
	
}

// RemoveByRefresh use the refresh token to delete the token information
func (tokenStorage *TokenStore) RemoveByRefresh(refresh string) (err error) {

}

func (tokenStorage *TokenStore) getData(basicID string) (ti oauth2.TokenInfo, err error) {
	
}

func (tokenStorage *TokenStore) getBasicID(cname, token string) (basicID string, err error) {

}

// GetByCode use the authorization code for token information data
func (tokenStorage *TokenStore) GetByCode(code string) (ti oauth2.TokenInfo, err error) {
	
}

// GetByAccess use the access token for token information data
func (tokenStorage *TokenStore) GetByAccess(access string) (ti oauth2.TokenInfo, err error) {

}

// GetByRefresh use the refresh token for token information data
func (tokenStorage *TokenStore) GetByRefresh(refresh string) (ti oauth2.TokenInfo, err error) {
	
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
