# syntax=docker/dockerfile:1

FROM golang:1.21

WORKDIR /wotlk

ENV NODE_VERSION=19.8.0
ENV NVM_DIR="/root/.nvm"

RUN apt-get update && \
        apt-get install -y protobuf-compiler

COPY go.mod go.sum ./
RUN go get -u google.golang.org/protobuf && go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
RUN curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.38.0/install.sh | bash
RUN . "$NVM_DIR/nvm.sh" && \
        nvm install ${NODE_VERSION} && \
        nvm use v${NODE_VERSION} && \
        nvm alias default v${NODE_VERSION}

ENV PATH="/root/.nvm/versions/node/v${NODE_VERSION}/bin/:${PATH}"

COPY gitconfig /etc/gitconfig
COPY . .

EXPOSE 8080/tcp
