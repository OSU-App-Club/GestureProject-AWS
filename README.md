# GestureProject-AWS

Run Go program

```bash
go mod tidy

go run .
```

# gRPC

Install protoc
```bash
sudo apt-get update
sudo apt-get install -y protobuf-compiler
sudo go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
sudo go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
```

Generate message .go files from .proto file
```bash
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    message/message.proto
```

Run server
```bash
go run server/main.go
```

Run client
```bash
go run client/main.go
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
