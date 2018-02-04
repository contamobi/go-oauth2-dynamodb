package dynamo

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
)

// Config dynamodb configuration parameters
type Config struct {
	SESSION  *session.Session
	TABLE    *TableConfig
	ENDPOINT string
}

type TableConfig struct {
	BasicCname   string
	AccessCName  string
	RefreshCName string
}

// NewConfig create dynamodb configuration
func NewConfig(region string, endpoint string, access_key string, secret string) (config *Config, err error) {
	newSession, err := session.NewSession(&aws.Config{
		Region:      aws.String(region),
		Credentials: credentials.NewStaticCredentials(access_key, secret, ""),
		Endpoint:    aws.String(endpoint),
	})
	if err != nil {
		return
	}
	config = &Config{
		SESSION: newSession,
		TABLE: &TableConfig{
			BasicCname:   "oauth2_basic",
			AccessCName:  "oauth2_access",
			RefreshCName: "oauth2_refresh",
		},
		ENDPOINT: endpoint,
	}
	return
}
