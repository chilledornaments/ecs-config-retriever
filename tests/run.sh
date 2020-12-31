#!/usr/bin/env bash

# Using flags
go run ./cmd/retriever/main.go -parameter=retriever-test -path=/tmp/param-not-encoded

go run ./cmd/retriever/main.go -parameter=retriever-test-encoded -encoded -path=/tmp/param-encoded

go run ./cmd/retriever/main.go -parameter=retriever-test-encoded-json -encoded -path=/tmp/param-encoded-json

# Using env vars

RETRIEVER_PARAMETER=retriever-test RETRIEVER_PATH=/tmp/param-no-env RETRIEVER_ENCODED=false go run ./cmd/retriever/main.go -from-env