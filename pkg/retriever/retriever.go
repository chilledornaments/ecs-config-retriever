package retriever

import (
	"os"

	"encoding/base64"
	"encoding/json"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	vault "github.com/hashicorp/vault/api"
	"github.com/sirupsen/logrus"
)

// GetParameterFromSSM retrieves the parameter from SSM
func GetParameterFromSSM(name string, encrypted, encoded bool, log *logrus.Logger) string {
	sess, err := session.NewSession(&aws.Config{Region: aws.String(os.Getenv("AWS_REGION"))})

	if err != nil {
		log.Fatalf("Error creating AWS Session: %s", err.Error())
	}

	log.Infof("Retrieving parameter '%s'", name)

	svc := ssm.New(sess)

	input := &ssm.GetParameterInput{
		Name:           aws.String(name),
		WithDecryption: aws.Bool(encrypted),
	}
	param, err := svc.GetParameter(input)

	if err != nil {
		log.Fatalf("Error retrieving parameter: %s", err.Error())
	}

	log.Info("Successfully retrieved parameter")

	if encoded {
		return decodeParameterValue(*param.Parameter.Value, log)
	}

	return *param.Parameter.Value

}

func GetSecretFromVault(path string, encoded bool, log *logrus.Logger, c *vault.Logical) map[string]string {

	secret, err := c.Read(path)

	if secret == nil {
		log.Fatalf("Secret returned from vault was nil")
	}

	if err != nil {
		log.Fatalf("Error reading secret from Vault path '%s': %s", path, err.Error())
	} else if len(secret.Warnings) > 0 {
		log.Fatalf("Errors returned from Vault: %v", secret.Warnings)
	}

	m := make(map[string]string)

	rv := secret.Data["data"]

	b, err := json.Marshal(rv)

	if err != nil {
		panic(err)
	}

	json.Unmarshal(b, &m)

	if encoded {
		newMap := make(map[string]string)

		for k, v := range m {
			newMap[k] = decodeParameterValue(v, log)
		}

		return newMap
	}

	// We return the entire map
	return m
}

// decodeParameterValue returns a base64-decoded string
func decodeParameterValue(value string, log *logrus.Logger) string {
	data, err := base64.StdEncoding.DecodeString(value)

	if err != nil {
		log.Fatalf("Error decoding parameter store value: %s", err.Error())
	}

	log.Info("Successfully decoded secret")

	return string(data)
}
