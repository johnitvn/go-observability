#!/bin/bash
set -e
lefthook install 
make install
make tidy 
docker create network ecoma-network || true
docker compose -f deploy/dev/docker-compose.infras.yaml up -d
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest