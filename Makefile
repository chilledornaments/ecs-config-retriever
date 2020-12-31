docker-build:
	docker build -t mitchya1/ecs-ssm-retriever:$(VERSION) .

push:
	docker push mitchya1/ecs-ssm-retriever:$(VERSION)

integration-tests:
	bash tests/run.sh

cleanup:
	rm /tmp/param-*