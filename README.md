# GestureProject-AWS

Run Go program

```bash
go run .
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
