FROM golang:1.20.0 AS build

# Set destination for COPY
WORKDIR /app

# Install protoc
RUN apt-get update
RUN apt-get install -y protobuf-compiler
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2

# Download Go modules
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy the source code. Note the slash at the end, as explained in
# https://docs.docker.com/engine/reference/builder/#copy
COPY . .

# Run protoc
RUN protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    message/message.proto

# Build
RUN go build -o /main ./*.go

## Deploy
# need a Docker image that will resolve error:  /lib/x86_64-linux-gnu/libm.so.6: version `GLIBC_2.29' not found (required by /main)
# like the golang image
FROM golang:1.20.0

WORKDIR /

COPY --from=build /main /main

ENTRYPOINT ["/main"]