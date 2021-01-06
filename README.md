# ecs-ssm-retriever

![Go Tests](https://github.com/mitchya1/ecs-ssm-retriever/workflows/Go%20Tests/badge.svg) ![Docker Tests](https://github.com/mitchya1/ecs-ssm-retriever/workflows/Docker%20Tests/badge.svg) ![CodeQL](https://github.com/mitchya1/ecs-ssm-retriever/workflows/CodeQL/badge.svg) [![Go Report Card](https://goreportcard.com/badge/github.com/mitchya1/ecs-ssm-retriever)](https://goreportcard.com/report/github.com/mitchya1/ecs-ssm-retriever) [![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=mitchya1_ecs-ssm-retriever&metric=alert_status)](https://sonarcloud.io/dashboard?id=mitchya1_ecs-ssm-retriever)

An init app used to retrieve configurations from SSM and write them to disk for use by a long running ECS container.

This is useful for writing a configuration file stored in SSM to a volume shared by ECS containers in a task. Currently, ECS doesn't allow you to mount a secret / configuration as a file like Kubernetes does. This tool works around that.

[Image on Docker Hub](https://hub.docker.com/r/mitchya1/ecs-ssm-retriever)

## Flags

`-parameter`: The name of the SSM Parameter Store parameter to retrieve

`-encoded`: Whether or not the parameter is base64 encoded. Default `false`

`-encrypted`: Whether or not the parameter is encrypted. Default: `false`

`-path`: The file to save the parameter to

`-from-env`: Specify this flag to tell `retriever` to get parameter info from the environment. Default: `false`. Conflicts with `-from-json`

`-from-json`: Specify this falg to tell `retriever` to get parameter info from a JSON passed as a string. Conflicts with `-from-env`

`-json`: JSON-as-a-string that specifies which parameters to retrieve. See the `JSON Argument` section for more information

## Env Vars

`AWS_REGION` - must be provided so an AWS session can be created. Set this to the region `retriever` is running in

`RETRIEVER_PARAMETER`: See `-parameter` flag

`RETRIEVER_PATH`: See `-path` flag

`RETRIEVER_ENCODED`: see `-encoded` flag

`RETRIEVER_ENCRYPTED`: see `-encrypted` flag

## JSON Argument

In order to retrieve multiple parameters, yYou can provide a JSON as a string to the `-json` argument.

JSON structure:

```json
{
    "parameters": [
        {
            "name": "some-parameter",
            "encoded": false,
            "encrypted": true,
            "path": "/init-out/some-app/some-parameter.yaml"
        },
        {
            "name": "some-other-parameter",
            "encoded": true,
            "encrypted": false,
            "path": "/init-out/some-other-app/some-other-parameter.json"
        }
    ]
}
```

## IAM Permissions

`retriever` needs minimal IAM permissions. This is the policy for the test suite user:

```json
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Action": "ssm:GetParameter",
            "Resource": "arn:aws:ssm:*:ACCOUNT_ID:parameter/retriever*"
        }
    ]
}
```

## Example Container ECS Definition

You must not change the `containerPath` for the `retriever` container, otherwise you'll receive permissions errors on file write.

```json
[
    {
        "command": ["cat", "/my-container/config/config.conf"],
        "cpu": 100,
        "essential": true,
        "environment": [
            {
                "name": "FOO",
                "value": "BAR"
            }
        ],  
        "mountPoints": [
            {
                "sourceVolume": "my-container-config",
                "containerPath": "/my-container/config"
            }
        ],
        "logConfiguration": {
            "logDriver": "awslogs",
            "options": {
                "awslogs-group": "my-container",
                "awslogs-region": "us-east-2",
                "awslogs-stream-prefix": "my-container"
            }
        },
        "volumesFrom": [],
        "image": "someimage:1.1.1",
        "memory": 100,
        "memoryReservation": 100,
        "name": "my-container",
        "dependsOn": [
            {
                "condition": "SUCCESS",
                "containerName": "my-container-init"
            }
        ],
        "privileged": false,
        "startTimeout": 60,
        "stopTimeout": 30
    },
    {
        "command": ["/retriever", "-parameter=retriever-test", "-path=/init-out/config.conf", "-encoded"],
        "cpu": 100,
        "essential": false,
        "environment": [
            {
                "name": "AWS_REGION",
                "value": "us-east-2"
            }
        ],  
        "mountPoints": [
            {
                "sourceVolume": "my-container-config",
                "containerPath": "/init-out"
            }
        ],
        "portMappings": [],
        "volumesFrom": [],
        "logConfiguration": {
            "logDriver": "awslogs",
            "options": {
                "awslogs-group": "my-container",
                "awslogs-region": "us-east-2",
                "awslogs-stream-prefix": "init"
            }
        },
        "image": "mitchya1/ecs-ssm-retriever:v0.2.2",
        "memory": 100,
        "memoryReservation": 50,
        "name": "my-container-init",
        "privileged": false,
        "startTimeout": 30,
        "stopTimeout": 60
    }
]
```

## Links

[Fargate shared volumes](https://docs.aws.amazon.com/AmazonECS/latest/developerguide/fargate-task-storage.html)

## Notes

The container initially runs as root so it can `chown` the `/init-out` directory. The command passed to the container is run as the non-privileged `retriever` user.
