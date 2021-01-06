FROM golang:1.15.6 AS builder

ENV GOOS=linux GOARCH=amd64 CGO_ENABLED=0

WORKDIR /go/src/github.com/mitchya1/ecs-ssm-retriever

COPY . .

RUN go build -o retriever ./cmd/retriever

FROM alpine:latest

RUN apk update && apk upgrade --no-cache && apk --no-cache add ca-certificates && rm -rf /var/cache/apk/*

COPY --from=builder /go/src/github.com/mitchya1/ecs-ssm-retriever/retriever /

RUN adduser --system --no-create-home --uid 121 retriever

RUN mkdir /init-out

VOLUME "/init-out"

RUN chown -R retriever /init-out

USER retriever