FROM migrate/migrate:latest AS golang-migrate

FROM golang:1.24.2 AS builder

RUN go install golang.org/x/tools/cmd/godoc@latest

FROM mcr.microsoft.com/devcontainers/go:1-1.23-bookworm

ARG USER=tech
ARG GROUP=$USER

RUN useradd -m $USER -s /bin/bash && \
    echo "$USER:$GROUP" | chpasswd && \
    echo "$USER ALL=(ALL) NOPASSWD:ALL" > /etc/sudoers.d/$USER && \
    chmod 0440 /etc/sudoers.d/$USER && \
    mkdir -p /home/$USER

COPY --from=golang-migrate /usr/local/bin/migrate /usr/local/bin/migrate
COPY --from=builder /go/bin/godoc /usr/local/bin/godoc

RUN apt-get update -y && \
    apt-get install -y --no-install-recommends clang-format && \
    apt-get clean && \
    rm -fr /var/lib/apt/lists/*

USER $USER
WORKDIR /workspaces/business-calender

RUN git config --global --add safe.directory /workspaces/business-calender
