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

Generate message .go files from .proto file within client and server directories

```bash
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    message/message.proto
```

Run server

```bash
cd server
go run .
```

Run client

```bash
cd client
go run .
```

# AWS Copilot CLI

-   For easy setup of ECS Fargate cluster + VPC & other stuff
-   https://aws.github.io/copilot-cli
-   https://aws.amazon.com/blogs/opensource/containerize-and-deploy-a-grpc-application-on-aws-fargate/

```bash
copilot init
- name app as gesture-project
- choose background service
- name service as data-processor
- select yes to create test environment

copilot app ls
copilot svc ls

copilot svc deploy --name server
copilot svc deploy --name client


copilot svc logs --follow --since 10m --name server
copilot svc logs --follow --since 10m --name client
```

# Docker

Build

```bash
sudo docker build -t gesture-project-aws-grpc-client .

sudo docker build -t gesture-project-aws-grpc-server .
```

OR Compose

```bash
sudo docker compose up --build
sudo docker compose up -d
```

Push to ECR

```bash
aws ecr get-login-password --region us-west-2 --profile aws-osuapp | docker login --username AWS --password-stdin 978103014270.dkr.ecr.us-west-2.amazonaws.com

docker tag gestureproject-dataprocessor:latest 978103014270.dkr.ecr.us-west-2.amazonaws.com/gestureproject-dataprocessor:latest

docker push 978103014270.dkr.ecr.us-west-2.amazonaws.com/gestureproject-dataprocessor:latest
```

# Note

If you run `copilot svc deploy` and get error during docker build process like the following:

```bash
# need a Docker image that will resolve error:  /lib/x86_64-linux-gnu/libm.so.6: version `GLIBC_2.29' not found (required by /main)
```

It may be because you are running on Mac (tested on Mac Mini M1), and it is solved by running the command in a Linux environment instead like GitHub Codespaces
