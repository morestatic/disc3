# syntax=docker/dockerfile:1

FROM golang:1.17-buster

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./
COPY env.docker ./.env
ADD archives ./archives
ADD cli ./cli
ADD commands ./commands
ADD documents ./documents
ADD indexes ./indexes
ADD ipfsutils ./ipfsutils
ADD progress ./progress
ADD server ./server
ADD types ./types
ADD utils ./utils


RUN go build -o /app/disc3

EXPOSE 4000

ENTRYPOINT [ "/app/disc3", "utils",  "repo-server" ]
