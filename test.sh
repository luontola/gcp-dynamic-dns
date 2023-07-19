#!/usr/bin/env bash
set -eu
cd src/app
go test ./...
#go test app/config app/gcloud
#go test app/ip -o ip.test.tmp
