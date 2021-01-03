docker-build:
	docker build -t mitchya1/ecs-ssm-retriever:$(VERSION) .

docker-push:
	docker push mitchya1/ecs-ssm-retriever:$(VERSION)

unit-tests:
	go test -v ./cmd/retriever/
	rm -rf /tmp/ci-*

integration-tests:
	bash tests/run.sh

cleanup:
	rm /tmp/param-*
	rm /tmp/binary-param-*

docker-cleanup:
	docker rmi mitchya1/ecs-ssm-retriever:$(VERSION)

docker-tests:
	docker-compose up