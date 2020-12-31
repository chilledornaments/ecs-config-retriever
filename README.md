# ecs-ssm-retriever

![Tests](https://github.com/mitchya1/ecs-ssm-retriever/workflows/Tests/badge.svg)

A dependant container used to retrieve configurations from SSM.

This is useful for writing a configuration file stored in SSM to a volume shared by ECS containers in a task. Currently, ECS doesn't allow you to mount a secret / configuration as a file like Kubernetes does. This tool works around that.

## Flags

`-parameter`: The name of the SSM Parameter Store parameter to retrieve

`-encoded`: Whether or not the parameter is base64 encoded. Default `false`

`-encrypted`: Whether or not the parameter is encrypted. Default: `false`

`-path`: The file to save the parameter to

`from-env`: Specify this flag to tell `retriever` to get its configuration from the environment. Default: `false`

## Env Vars

`AWS_REGION` - must be provided so an AWS session can be created. Set this to the region `retriever` is running in

`RETRIEVER_PARAMETER`: See `-parameter` flag

`RETRIEVER_PATH`: See `-path` flag

`RETRIEVER_ENCODED`: see `-encoded` flag

`RETRIEVER_ENCRYPTED`: see `-encrypted` flag

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
        "command": ["cat", "/my-container/config/"],
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
        "command": ["/retriever", "-parameter=retriever-test", "-path=/init-out/param-not", "-encoded"],
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
        "image": "mitchya1/ecs-ssm-retriever:v0.1.0",
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


## Future State

Retrieve multiple parameters as once