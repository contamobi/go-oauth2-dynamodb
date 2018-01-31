package dynamo

import (
	"encoding/json"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"gopkg.in/mgo.v2/bson"
	//"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/contamobi/oauth2"
)

func NewTokenStore(config *Config) (tokenStore *TokenStore) {
	session := config.SESSION
	svc := dynamodb.New(session)
	return &TokenStore{
		config:  config,
		session: svc,
	}
}

type TokenStore struct {
	config  *Config
	session *dynamodb.DynamoDB
}

type tokenData struct {
	ID        string
	BasicID   string
	ExpiredAt time.Time
}

type basicData struct {
	ID        string
	Data      []byte
	ExpiredAt time.Time
}

// Create and store the new token information
func (tokenStorage *TokenStore) Create(info oauth2.TokenInfo) (err error) {
	jv, err := json.Marshal(info)
	if err != nil {
		return
	}

	if code := info.GetCode(); code != "" {
		err = CreateWithAuthorizationCode(tokenStorage, info)
		return
	}
	if refresh := info.GetRefresh(); refresh != "" {
		err = CreateWithRefreshToken(tokenStorage, info)
	} else {
		err = CreateWithAccessToken(tokenStorage, info)
	}
	return
}

func CreateWithAuthorizationCode(tokenStorage *TokenStore, info oauth2.TokenInfo, id ...string) (err error) {
	code := info.GetCode()

	if id != nil {
		code = id[0]
	}
	data, err := json.Marshal(info)
	if err != nil {
		return
	}
	expiredAt := info.GetCodeCreateAt().Add(info.GetCodeExpiresIn()).String()
	rExpiredAt := expiredAt
	if refresh := info.GetRefresh(); refresh != "" {
		rexp := info.GetRefreshCreateAt().Add(info.GetRefreshExpiresIn())
		if info.GetCodeCreateAt().Add(info.GetCodeExpiresIn()).Second() > rexp.Second() {
			expiredAt = rexp.String()
		}
		rExpiredAt = rexp.String()
	}
	params := &dynamodb.PutItemInput{
		TableName: aws.String(tokenStorage.config.TABLE.TxnCName),
		Item: map[string]*dynamodb.AttributeValue{
			"ID": &dynamodb.AttributeValue{
				S: aws.String(code),
			},
			"Data": &dynamodb.AttributeValue{
				B: data,
			},
			"ExpiredAt": &dynamodb.AttributeValue{
				S: &rExpiredAt,
			},
		},
	}
	_, err = tokenStorage.session.PutItem(params)
	return
}

func CreateWithAccessToken(tokenStorage *TokenStore, info oauth2.TokenInfo, id ...string) (err error) {

	if id == nil {
		id[0] = bson.NewObjectId().Hex()
	}
	err = CreateWithAuthorizationCode(tokenStorage, info, id[0])
	if err != nil {
		return
	}
	expiredAt := info.GetAccessCreateAt().Add(info.GetAccessExpiresIn()).String()
	token := tokenData{
		BasicID:   id[0],
		ExpiredAt: info.GetAccessCreateAt().Add(info.GetAccessExpiresIn()),
	}
	tokenByte, err := json.Marshal(token)
	if err != nil {
		return
	}
	accessParams := &dynamodb.PutItemInput{
		TableName: aws.String(tokenStorage.config.TABLE.AccessCName),
		Item: map[string]*dynamodb.AttributeValue{
			"ID": &dynamodb.AttributeValue{
				S: aws.String(info.GetAccess()),
			},
			"Data": &dynamodb.AttributeValue{
				B: tokenByte,
			},
			"ExpiredAt": &dynamodb.AttributeValue{
				S: &expiredAt,
			},
		},
	}
	_, err = tokenStorage.session.PutItem(accessParams)
	return
}

func CreateWithRefreshToken(tokenStorage *TokenStore, info oauth2.TokenInfo) (err error) {
	id := bson.NewObjectId().Hex()
	err = CreateWithAccessToken(tokenStorage, info, id)
	if err != nil {
		return
	}
	expiredAt := info.GetRefreshCreateAt().Add(info.GetRefreshExpiresIn()).String()
	token := tokenData{
		BasicID:   id,
		ExpiredAt: info.GetRefreshCreateAt().Add(info.GetRefreshExpiresIn()),
	}
	tokenByte, err := json.Marshal(token)
	if err != nil {
		return
	}
	accessParams := &dynamodb.PutItemInput{
		TableName: aws.String(tokenStorage.config.TABLE.RefreshCName),
		Item: map[string]*dynamodb.AttributeValue{
			"ID": &dynamodb.AttributeValue{
				S: aws.String(info.GetRefresh()),
			},
			"Data": &dynamodb.AttributeValue{
				B: tokenByte,
			},
			"ExpiredAt": &dynamodb.AttributeValue{
				S: &expiredAt,
			},
		},
	}
	_, err = tokenStorage.session.PutItem(accessParams)
	return
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
