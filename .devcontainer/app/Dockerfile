FROM migrate/migrate:latest AS golang-migrate

FROM golang:1.24.2 AS builder

RUN go install golang.org/x/tools/cmd/godoc@latest

FROM mcr.microsoft.com/devcontainers/go:1-1.23-bookworm

COPY --from=golang-migrate /usr/local/bin/migrate /usr/local/bin/migrate
COPY --from=builder /go/bin/godoc /usr/local/bin/godoc

RUN apt-get update -y && \
    apt-get install -y --no-install-recommends clang-format && \
    apt-get clean && \
    rm -fr /var/lib/apt/lists/*
