docker-build:
	docker build -t mitchya1/ecs-ssm-retriever:$(VERSION) .

push:
	docker push mitchya1/ecs-ssm-retriever:$(VERSION)

integration-tests:
	bash tests/run.sh

cleanup:
	rm /tmp/param-*
	rm /tmp/binary-param-*

docker-cleanup:
	docker rmi mitchya1/ecs-ssm-retriever:$(VERSION)