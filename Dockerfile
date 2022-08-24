# syntax=docker/dockerfile:1
#build stage

FROM golang:alpine AS builder
FROM golang:1.18
RUN apk add --no-cache git

WORKDIR /go/src/app
WORKDIR /wotlk
COPY . .
COPY . .
RUN go get -d -v ./...

RUN go build -o /go/bin/app -v ./...
RUN apt-get update

RUN apt-get install -y protobuf-compiler
#final stage
RUN go get -u google.golang.org/protobuf
FROM alpine:latest
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
RUN apk --no-cache add ca-certificates

COPY --from=builder /go/bin/app /app
RUN curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.38.0/install.sh | bash
ENTRYPOINT /app

LABEL Name=wotlk Version=0.0.1
ENV NODE_VERSION=14.18.3
EXPOSE 8080
ENV NVM_DIR="/root/.nvm"
RUN . "$NVM_DIR/nvm.sh" && nvm install ${NODE_VERSION}
RUN . "$NVM_DIR/nvm.sh" && nvm use v${NODE_VERSION}
RUN . "$NVM_DIR/nvm.sh" && nvm alias default v${NODE_VERSION}
ENV PATH="/root/.nvm/versions/node/v${NODE_VERSION}/bin/:${PATH}"

EXPOSE 8080/tcp