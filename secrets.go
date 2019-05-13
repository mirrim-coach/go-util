package goutils

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"os"
	"regexp"
	"strings"
)

var stdKeyConversions = map[string]map[string]string{
	"SECRETS_DB": map[string]string{
		"username": "DB_USER",
		"password": "DB_PASS",
		"dbname":   "DB_NAME",
		"port":     "DB_PORT",
		"host":     "DB_HOST",
		"ssl":      "DB_SSL",
	},
}

func convertKey(secretEnv string, key string) string {
	if m, ok := stdKeyConversions[secretEnv]; ok {
		if v, ok := m[key]; ok {
			key = v
		}
	}

	return key
}

func getSecret(secretEnv string, secretName string) {
	//Create a Secrets Manager client
	svc := secretsmanager.New(GetAWSSession())
	input := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(secretName),
		VersionStage: aws.String("AWSCURRENT"), // VersionStage defaults to AWSCURRENT if unspecified
	}

	// In this sample we only handle the specific exceptions for the 'GetSecretValue' API.
	// See https://docs.aws.amazon.com/secretsmanager/latest/apireference/API_GetSecretValue.html

	result, err := svc.GetSecretValue(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case secretsmanager.ErrCodeDecryptionFailure:
				// Secrets Manager can't decrypt the protected secret text using the provided KMS key.
				Logger().Error(secretsmanager.ErrCodeDecryptionFailure, aerr.Error())

			case secretsmanager.ErrCodeInternalServiceError:
				// An error occurred on the server side.
				Logger().Error(secretsmanager.ErrCodeInternalServiceError, aerr.Error())

			case secretsmanager.ErrCodeInvalidParameterException:
				// You provided an invalid value for a parameter.
				Logger().Error(secretsmanager.ErrCodeInvalidParameterException, aerr.Error())

			case secretsmanager.ErrCodeInvalidRequestException:
				// You provided a parameter value that is not valid for the current state of the resource.
				Logger().Error(secretsmanager.ErrCodeInvalidRequestException, aerr.Error())

			case secretsmanager.ErrCodeResourceNotFoundException:
				// We can't find the resource that you asked for.
				Logger().Error(secretsmanager.ErrCodeResourceNotFoundException, aerr.Error())
			default:
				Logger().Error(aerr.Error())
			}
		} else {
			Logger().Error(err.Error())
		}
		return
	}

	// Decrypts secret using the associated KMS CMK.
	// Depending on whether the secret is a string or binary, one of these fields will be populated.
	secretString := ""
	if result.SecretString != nil {
		secretString = *result.SecretString
	} else {
		decodedBinarySecretBytes := make([]byte, base64.StdEncoding.DecodedLen(len(result.SecretBinary)))
		len, err := base64.StdEncoding.Decode(decodedBinarySecretBytes, result.SecretBinary)
		if err != nil {
			Logger().Error("Base64 Decode Error:", err)
			return
		}
		secretString = string(decodedBinarySecretBytes[:len])
	}

	variables := map[string]interface{}{}
	err = json.NewDecoder(bytes.NewBuffer([]byte(secretString))).Decode(&variables)

	if err != nil {
		panic(err)
	}

	Logger().Infof("Setting ENV Variables from %s", secretEnv)

	for k, v := range variables {
		k = convertKey(secretEnv, k)
		Logger().Infof("- Setting %s", k)
		os.Setenv(k, fmt.Sprintf("%v", v))
	}

	// Your code goes here.
}

//ConfigureEnvironmentFromSecrets sets env variables from secrets
func ConfigureEnvironmentFromSecrets() {
	envs := os.Environ()
	re := regexp.MustCompile(`(?m)^SECRETS_.+`)

	for _, v := range envs {
		split := strings.SplitN(v, "=", 2)
		if re.Match([]byte(split[0])) {
			Logger().Infof("Loading `%s` from secrets manager!", split[1])
			getSecret(split[0], split[1])
		}
	}
}
