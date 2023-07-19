#!/usr/bin/env bash
set -eux
docker-compose build --force-rm
docker-compose run --rm app "$@"
