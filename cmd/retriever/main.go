package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/mitchya1/ecs-file-retriever/pkg/retriever"

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

func main() {
	log := logrus.New()

	log.SetFormatter(&logrus.JSONFormatter{})
	log.SetOutput(os.Stdout)

	log.Info("Starting ECS SSM Bootstrapper")

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
		j := parseJSONArgument(log)
		for _, p := range j.Parameters {
			v := retriever.GetParameterFromSSM(p.Name, p.Encryped, p.Encoded, log)
			createDirectory(p.Path, log)
			writeSecretToFile(v, p.Path, log)
		}
		// return so we don't continue processesing
		return
	}

	if fromVault {
		// TODO add support for AgentAddress
		// TODO make this more configurable
		vc := vault.Config{
			Address:    os.Getenv("VAULT_ADDR"),
			MaxRetries: 2,
			Timeout:    4 * time.Second,
		}

		v, err := vault.NewClient(&vc)

		if err != nil {
			log.Fatalf("Error creating Vault client: %s", err.Error())
		}

		// TODO add support for STS
		v.SetToken(os.Getenv("VAULT_TOKEN"))

		c := v.Logical()
		m := retriever.GetSecretFromVault(vaultPath, parameterIsEncoded, log, c)

		s := new(bytes.Buffer)

		for k, v := range m {
			fmt.Fprintf(s, "%s = %s\n", k, v)
		}

		writeSecretToFile(s.String(), filePath, log)
		return
	}

	if fromEnv {
		getValuesFromEnv(log)
	}

	v := retriever.GetParameterFromSSM(parameterName, parameterIsEncrypted, parameterIsEncoded, log)
	createDirectory(filePath, log)
	writeSecretToFile(v, filePath, log)
}
