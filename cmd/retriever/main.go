package main

import (
	"flag"
	"os"
	"strconv"

	"github.com/mitchya1/ecs-ssm-retriever/pkg/retriever"
	"github.com/sirupsen/logrus"
)

var (
	fromEnv              bool
	parameterIsEncoded   bool
	parameterIsEncrypted bool
	parameterName        string
	filePath             string
)

func main() {
	log := logrus.New()

	log.SetFormatter(&logrus.JSONFormatter{})
	log.SetOutput(os.Stdout)

	log.Info("Starting ECS SSM Bootstrapper")

	flag.BoolVar(&fromEnv, "from-env", false, "Retrieve settings from env")
	flag.StringVar(&parameterName, "parameter", os.Getenv("PARAMETER"), "Name of parameter to retrieve")
	flag.BoolVar(&parameterIsEncoded, "encoded", false, "Decides whether or not the parameter will be base64 decoded prior to writing to file")
	flag.BoolVar(&parameterIsEncrypted, "encrypted", false, "If the SSM parameter is encrypted, provide this argument")
	flag.StringVar(&filePath, "path", "", "Path to save retrieved parameter to")

	flag.Parse()

	if fromEnv {
		getValuesFromEnv(log)
	}

	v := retriever.GetParameterFromSSM(parameterName, parameterIsEncrypted, parameterIsEncoded, log)

	writeSecretToFile(v, filePath, log)
}

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

	log.Infof("Successfully wrote paramater to '%s'", path)
}

// getValuesFromEnv sets variable values by getting them from the environment, not CLI args
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
