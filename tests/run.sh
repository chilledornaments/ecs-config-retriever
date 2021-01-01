#!/usr/bin/env bash

# Using env vars

RETRIEVER_PARAMETER=retriever-test RETRIEVER_PATH=/tmp/param-not-encoded-env RETRIEVER_ENCODED=false go run ./cmd/retriever/main.go -from-env

# Using binary

./retriever -parameter=retriever-test -path=/tmp/binary-param-not-encoded

./retriever -parameter=retriever-test-encoded -encoded -path=/tmp/binary-param-encoded

./retriever -parameter=retriever-test-encoded-json -encoded -path=/tmp/binary-param-encoded-json

RETRIEVER_PARAMETER=retriever-test RETRIEVER_PATH=/tmp/binary-param-not-encoded-env RETRIEVER_ENCODED=false ./retriever -from-env

./retriever -from-json -json '{"parameters": [{"name": "retriever-test", "encrypted": false, "encoded": false, "path": "/tmp/param-json-0"}]}'

./retriever -from-json -json '{"parameters": [{"name": "retriever-test-encoded", "encrypted": false, "encoded": true, "path": "/tmp/param-json-1"}]}'

# This check should fail
./retriever -from-json -from-env || exit 0

# This check should fail
./retriever -from-json || exit 0