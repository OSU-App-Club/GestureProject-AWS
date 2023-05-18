# GestureProject-AWS

Run Go program

```bash
go mod tidy

go run .
```

# AWS Copilot CLI

-   For easy setup of ECS Fargate cluster + VPC & other stuff
-   https://aws.github.io/copilot-cli

```bash
copilot init
- name app as gesture-project
- choose background service
- name service as data-processor
- select yes to create test environment

copilot app ls

copilot svc deploy

copilot svc logs --follow --since 1h
```

# Docker

Build

```bash
docker build -t gestureproject-dataprocessor .
```

OR Compose

```bash
docker compose up --build
```

Push to ECR

```bash
aws ecr get-login-password --region us-west-2 --profile aws-osuapp | docker login --username AWS --password-stdin 978103014270.dkr.ecr.us-west-2.amazonaws.com

docker tag gestureproject-dataprocessor:latest 978103014270.dkr.ecr.us-west-2.amazonaws.com/gestureproject-dataprocessor:latest

docker push 978103014270.dkr.ecr.us-west-2.amazonaws.com/gestureproject-dataprocessor:latest
```
