FROM golang:1.16.1 AS builder

ENV GOOS=linux GOARCH=amd64 CGO_ENABLED=0

WORKDIR /go/src/github.com/mitchya1/ecs-config-retriever

COPY . .

RUN go build -o retriever ./cmd/retriever

FROM alpine:latest

RUN apk update && apk upgrade --no-cache && apk --no-cache add ca-certificates su-exec && rm -rf /var/cache/apk/*

COPY --from=builder /go/src/github.com/mitchya1/ecs-config-retriever/retriever /

ADD docker-entrypoint.sh /

RUN chmod +x /docker-entrypoint.sh

RUN adduser --system --no-create-home --uid 121 retriever

RUN mkdir /init-out

VOLUME "/init-out"

RUN chown -R retriever /init-out

ENTRYPOINT [ "/docker-entrypoint.sh" ]

# Run as root so entrypoint can chown the /init-out dir then su-exec as retriever
USER root