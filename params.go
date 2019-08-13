package main

import (
	"fmt"
	"os"
  "encoding/json"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
)

func ConfigureEnvironmentFromParameterStore(parameterStoreName string, region string) {
	fmt.Println("Getting parameter:",parameterStoreName)
	sess, err := session.NewSessionWithOptions(session.Options{
		Config:            aws.Config{Region: aws.String(region)},
		SharedConfigState: session.SharedConfigEnable,
	})
	if err != nil {
		panic(err)
	}

	ssmsvc := ssm.New(sess, aws.NewConfig().WithRegion(region))
	withDecryption := true
	param, err := ssmsvc.GetParameter(&ssm.GetParameterInput{
		Name:           &parameterStoreName,
		WithDecryption: &withDecryption,
	})
	parametersJsonByte := []byte(*param.Parameter.Value)
	fmt.Println("Parsing json")
	var paramsInterface interface{}
	json.Unmarshal(parametersJsonByte, &paramsInterface)
	parameters := paramsInterface.(map[string]interface{})

	fmt.Println("Setting parameters to environmnet variables")

	for key, value := range parameters {
		keyString := fmt.Sprintf("%v", key)
		valueString := fmt.Sprintf("%v", value)
		os.Setenv(keyString, valueString)
	}
}
