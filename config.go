package dynamo

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/aws/credentials"
)

// Config dynamodb configuration parameters
type Config struct {
	SESSION *session.Session
	TABLE  string
	ENDPOINT string
}

// NewConfig create dynamodb configuration
func NewConfig(region string, endpoint string, access_key string, secret string, table string) (config *Config, err error) {
	newSession, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
		Credentials: credentials.NewStaticCredentials(access_key, secret, ""),
	})
	if err != nil {
		return
	}
	config = &Config{
		SESSION: newSession,
		TABLE: table,
		ENDPOINT: endpoint,
	}
	return 
}
