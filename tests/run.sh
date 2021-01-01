#!/usr/bin/env bash

# Using flags
go run ./cmd/retriever/main.go -parameter=retriever-test -path=/tmp/param-not-encoded

go run ./cmd/retriever/main.go -parameter=retriever-test-encoded -encoded -path=/tmp/param-encoded

go run ./cmd/retriever/main.go -parameter=retriever-test-encoded-json -encoded -path=/tmp/param-encoded-json

# Using env vars

RETRIEVER_PARAMETER=retriever-test RETRIEVER_PATH=/tmp/param-not-encoded-env RETRIEVER_ENCODED=false go run ./cmd/retriever/main.go -from-env

# Using binary

./retriever -parameter=retriever-test -path=/tmp/binary-param-not-encoded

./retriever -parameter=retriever-test-encoded -encoded -path=/tmp/binary-param-encoded

./retriever -parameter=retriever-test-encoded-json -encoded -path=/tmp/binary-param-encoded-json

RETRIEVER_PARAMETER=retriever-test RETRIEVER_PATH=/tmp/binary-param-not-encoded-env RETRIEVER_ENCODED=false ./retriever -from-env
