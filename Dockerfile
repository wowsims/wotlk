# syntax=docker/dockerfile:1

FROM golang:1.21

WORKDIR /wotlk
COPY . .
COPY gitconfig /etc/gitconfig

RUN apt-get update && apt-get install -y \
    protobuf-compiler \
    curl

# Install protobuf
RUN go get -u google.golang.org/protobuf
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest

# Install Node.js 19.x using NodeSource
RUN curl -fsSL https://deb.nodesource.com/setup_19.x | bash - && \
    apt-get install -y nodejs

EXPOSE 8080/tcp
