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
	docker rmi --force mitchya1/ecs-ssm-retriever:$(VERSION)

docker-tests:
	@echo "ACCESS_KEY=${ACCESS_KEY}" > .env
	@echo "SECRET_KEY=${SECRET_KEY}" >> .env
	docker-compose up | tee /tmp/ci-compose-out
	grep "with code 1" /tmp/ci-compose-out && exit 1 || exit 0