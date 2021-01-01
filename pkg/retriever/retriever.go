package retriever

import (
	"os"

	"encoding/base64"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
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

// decodeParameterValue returns a base64-decoded string
func decodeParameterValue(value string, log *logrus.Logger) string {
	data, err := base64.StdEncoding.DecodeString(value)

	if err != nil {
		log.Fatalf("Error decoding parameter store value: %s", err.Error())
	}

	log.Info("Successfully decoded secret")

	return string(data)
}
