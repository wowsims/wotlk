# syntax=docker/dockerfile:1
FROM golang:1.21

WORKDIR /wotlk
COPY . .
COPY gitconfig /etc/gitconfig

RUN apt-get update
RUN apt-get install -y protobuf-compiler
RUN go get -u google.golang.org/protobuf
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest

ENV NODE_VERSION=18.18.0
ENV PATH="/opt/node-v${NODE_VERSION}-linux-x64/bin:${PATH}"
RUN curl https://nodejs.org/dist/v${NODE_VERSION}/node-v${NODE_VERSION}-linux-x64.tar.gz |tar xzf - -C /opt/

EXPOSE 8080/tcp
