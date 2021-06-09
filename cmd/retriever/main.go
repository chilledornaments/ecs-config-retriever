package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/mitchya1/ecs-config-retriever/pkg/retriever"

	vault "github.com/hashicorp/vault/api"
	"github.com/sirupsen/logrus"
)

var (
	fromEnv              bool
	fromJSON             bool
	fromVault            bool
	vaultPath            string
	vaultUseSTS          bool
	parameterIsEncoded   bool
	parameterIsEncrypted bool
	parameterName        string
	filePath             string
	jsonSettings         string
)

// JSONArgument holds our ParameterSettings passed with the -from-json flag
type JSONArgument struct {
	Parameters []ParameterSetting `json:"parameters"`
}

// ParameterSetting contains information about a Parameter Store parameter and where to write it out
type ParameterSetting struct {
	Name     string `json:"name"`
	Encryped bool   `json:"encrypted"`
	Encoded  bool   `json:"encoded"`
	Path     string `json:"path"`
}

func newAWSConfig(region string) aws.Config {
	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithRegion(region),
	)

	if err != nil {
		log.Fatalf("Failed to load SDK configuration, %v", err)
	}

	return cfg
}

func newSSMClient(cfg aws.Config) *ssm.Client {
	return ssm.NewFromConfig(cfg)
}

func ssmHandler(log *logrus.Logger) {

	awsConfig := newAWSConfig(os.Getenv("AWS_REGION"))

	ssmClient := newSSMClient(awsConfig)

	v, e := retriever.GetParameterFromSSM(context.Background(), ssmClient, log, parameterName, parameterIsEncrypted, parameterIsEncoded)

	if e != nil {
		// GetParameterFromSSM already logs the error
		os.Exit(1)
	}

	if e = createDirectory(filePath, log); e != nil {
		os.Exit(1)
	}
	if e = writeValueToFile(v, filePath, log); e != nil {
		os.Exit(1)
	}

}

func ssmJSONHandler(log *logrus.Logger, j JSONArgument) {
	awsConfig := newAWSConfig(os.Getenv("AWS_REGION"))

	ssmClient := newSSMClient(awsConfig)

	for _, p := range j.Parameters {
		v, e := retriever.GetParameterFromSSM(context.Background(), ssmClient, log, p.Name, p.Encryped, p.Encoded)

		if e != nil {
			// GetParameterFromSSM already logs the error
			os.Exit(1)
		}

		e = createDirectory(p.Path, log)

		if e != nil {
			// writeValueToFile already logs, so we won't do that again here
			os.Exit(1)
		}

		e = writeValueToFile(v, p.Path, log)

		if e != nil {
			// writeValueToFile already logs, so we won't do that again here
			os.Exit(1)
		}
	}
}

func vaultHandler(log *logrus.Logger) error {
	vc := vault.Config{
		Address:    os.Getenv("VAULT_ADDR"),
		MaxRetries: 2,
		Timeout:    4 * time.Second,
	}

	v, err := vault.NewClient(&vc)

	if err != nil {
		log.Errorf("Error creating Vault client: %s", err.Error())
		return err
	}

	// TODO add support for STS
	v.SetToken(os.Getenv("VAULT_TOKEN"))

	c := v.Logical()
	m := retriever.GetSecretFromVault(vaultPath, parameterIsEncoded, log, c)

	s := new(bytes.Buffer)

	for k, v := range m {
		fmt.Fprintf(s, "%s = %s\n", k, v)
	}

	if err = writeValueToFile(s.String(), filePath, log); err != nil {
		return err
	}

	return nil
}

func main() {
	log := logrus.New()

	log.SetFormatter(&logrus.JSONFormatter{})
	log.SetOutput(os.Stdout)

	log.Info("Starting ECS File Retriever")

	flag.BoolVar(&fromEnv, "from-env", false, "Retrieve settings from env")
	flag.BoolVar(&fromVault, "from-vault", false, "Retrieve settings from Hashi Vault")
	// If you enable the kv engine in the default 'kv' path, your path will look something like kv/data/foo/test
	// Note the `data/` path
	flag.StringVar(&vaultPath, "vault-path", "", "Path to Vault secret")
	flag.BoolVar(&vaultUseSTS, "vault-use-sts", false, "If Retriever can access Vault use an IAM role, set this flag")
	flag.StringVar(&parameterName, "parameter", "", "Name of parameter to retrieve")
	flag.BoolVar(&parameterIsEncoded, "encoded", false, "Decides whether or not the parameter will be base64 decoded prior to writing to file")
	flag.BoolVar(&parameterIsEncrypted, "encrypted", false, "If the SSM parameter is encrypted, provide this argument")
	flag.StringVar(&filePath, "path", "", "Path to save retrieved parameter to")
	flag.BoolVar(&fromJSON, "from-json", false, "Provide a JSON object of parameters to retrieve. Allows retrieving multiple parameters")
	flag.StringVar(&jsonSettings, "json", "", "JSON object of parameters to retrieve")

	flag.Parse()

	// Ensure that conflicting or incomplete arguments have not been provided
	verifyFlags(log)

	if fromJSON {
		j, e := parseJSONArgument(log)

		if e != nil {
			os.Exit(1)
		}

		ssmJSONHandler(log, j)
		// Return so we don't continue processing
		return
	}

	if fromVault {
		vaultHandler(log)
		// Return so we don't continue processing
		return
	}

	if fromEnv {
		getValuesFromEnv(log)
	}

	// Fall back to retrieving a single SSM parameter
	// This handles either env vars or args
	ssmHandler(log)
}
