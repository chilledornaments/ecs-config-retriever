package main

import (
	"os"
	"testing"

	"github.com/sirupsen/logrus"
	"gotest.tools/assert"
)

func setupTestingLogger() *logrus.Logger {
	log := logrus.New()

	log.SetFormatter(&logrus.JSONFormatter{})
	log.SetOutput(os.Stdout)
	return log
}
func TestParseJSONArgument(t *testing.T) {
	jsonSettings = `{"parameters":[
		{
			"name": "ci",
			"encrypted": false,
			"encoded": false,
			"path": "/hello/world.txt"
		}
	]}`

	l := setupTestingLogger()

	j, e := parseJSONArgument(l)

	assert.Equal(t, e, nil)

	if j.Parameters[0].Name != "ci" {
		l.Warnf("Expected parameter name to be 'ci' but it is '%s'", j.Parameters[0].Name)
		t.Fail()
	}

	if j.Parameters[0].Encoded {
		l.Warn("Expected parameter encoded to be 'false' but it is 'true'")
		t.Fail()
	}
}

func TestCreateDirectory(t *testing.T) {
	l := setupTestingLogger()
	createDirectory("/tmp/ci-test-dir/file.txt", l)

	f, err := os.Stat("/tmp/ci-test-dir")
	if err != nil {
		l.Warnf("Checking for /tmp/ci-test-dir raised an error: %s", err.Error())
		t.Fail()
	}
	if !f.IsDir() {
		l.Warn("/tmp/ci-test-dir is not a directory")
		t.Fail()
	}
}

func TestWriteSecretToFile(t *testing.T) {
	l := setupTestingLogger()
	e := writeValueToFile("hello world", "/tmp/file.txt", l)

	assert.Equal(t, e, nil)

}
