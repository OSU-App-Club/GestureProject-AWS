FROM golang:1.20.0 AS build

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy the source code. Note the slash at the end, as explained in
# https://docs.docker.com/engine/reference/builder/#copy
COPY . .

# Build
RUN go build -o /main ./*.go

## Deploy
# need a Docker image that will resolve error:  /lib/x86_64-linux-gnu/libm.so.6: version `GLIBC_2.29' not found (required by /main)
# like the golang image
FROM golang:1.20.0

WORKDIR /

COPY --from=build /main /main

ENTRYPOINT ["/main"]