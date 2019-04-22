package goutils

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

func GetAWSSession() *session.Session {
	s, err := session.NewSession(&aws.Config{Region: aws.String(GetEnvVariable("AWS_REGION", "us-east-1"))})
	if err != nil {
		Logger().Fatal(err)
	}
	return s
}
