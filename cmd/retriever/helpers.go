package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strconv"

	"github.com/sirupsen/logrus"
)

// writeSecretToFile writes the retrieved parameter value (value) to the specified path (path) for use between containers
func writeSecretToFile(value, path string, log *logrus.Logger) {
	f, err := os.Create(path)
	if err != nil {
		log.Fatalf("Error creating file: %s", err.Error())
	}

	defer f.Close()

	_, err = f.WriteString(value)

	if err != nil {
		log.Fatalf("Error writing parameter to file: %s", err.Error())
	}

	f.Sync()

	log.Infof("Successfully wrote parameter to '%s'", path)
}

// createDirectory creates the directory for the parameter out file to be stored in
// This is useful if you need to write files into subdirectories of your volume
// For instance, you can mount one volume onto /init-out/app-a/config and another onto /init-out/app-b/config
// Then mount these onto separate app containers
func createDirectory(path string, log *logrus.Logger) {
	fp := filepath.Dir(path)

	info, err := os.Stat(fp)

	if err != nil {
		log.Infof("Path '%s' does not exist. Attempting to create so we can store file", fp)
		err = os.MkdirAll(fp, 0775)
		if err != nil {
			log.Fatalf("Error creating directory structure '%s': %s", fp, err.Error())
		}
		log.Infof("Successfully created directory '%s'", fp)
	} else {
		if !info.IsDir() {
			log.Fatalf("'%s' is a file - unable to create directory in its place", fp)
		}
	}

}

// getValuesFromEnv retrieves configuration from env vars
func getValuesFromEnv(log *logrus.Logger) {
	var err error

	parameterName = os.Getenv("RETRIEVER_PARAMETER")
	filePath = os.Getenv("RETRIEVER_PATH")

	if os.Getenv("RETRIEVER_ENCODED") != "" {
		parameterIsEncoded, err = strconv.ParseBool(os.Getenv("RETRIEVER_ENCODED"))
		if err != nil {
			log.Fatalf("Unable to convert '%s' to bool", os.Getenv("RETRIEVER_ENCODED"))
		}
		log.Infof("Setting parameterIsEncoded to '%t'", parameterIsEncoded)
	} else {
		log.Info("RETRIEVER_ENCODED env var not set, defaulting to false")
	}

	if os.Getenv("RETRIEVER_ENCRYPTED") != "" {
		parameterIsEncrypted, err = strconv.ParseBool(os.Getenv("RETRIEVER_ENCRYPTED"))
		if err != nil {
			log.Fatalf("Unable to convert '%s' to bool", os.Getenv("RETRIEVER_ENCRYPTED"))
		}
		log.Infof("Setting parameterIsEncrypted to '%t'", parameterIsEncrypted)
	} else {
		log.Info("RETRIEVER_ENCRYPTED env var not set, defaulting to false")
	}
}

// verifyFlags ensures no flag conflicts or major issues
func verifyFlags(log *logrus.Logger) {

	if fromVault && parameterName != "" {
		log.Fatal("Cannot specify both -from-vault and -parameter")
	}

	if fromVault && os.Getenv("VAULT_ADDR") == "" {
		log.Fatal("-from-vault specified but VAULT_ADDR env var is empty")
	}

	// TODO don't require VAULT_TOKEN because we could use STS
	if fromVault && os.Getenv("VAULT_TOKEN") == "" && !vaultUseSTS {
		log.Fatal("-from-vault specified (not using STS) but VAULT_TOKEN env var is empty")
	}

	if fromEnv && fromJSON {
		log.Fatal("Cannot set -from-env and -from-json")
	}

	if fromJSON && jsonSettings == "" {
		log.Fatal("-from-json specified but no value provided for -json")
	}
}

// parseJSONArgument parses the -json argument into a struct
func parseJSONArgument(log *logrus.Logger) JSONArgument {
	j := &JSONArgument{}

	err := json.Unmarshal([]byte(jsonSettings), j)

	if err != nil {
		log.Fatalf("Unable to unmarshal -json argument into JSONArgument struct: %s", err.Error())
	}

	return *j
}
