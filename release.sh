#!/usr/bin/env bash
set -eu
: ${1:? Usage: $0 RELEASE_VERSION}
SCRIPTS=$(dirname "$0")

RELEASE_VERSION="$1"
if [[ ! "$RELEASE_VERSION" =~ ^[0-9]+\.[0-9]+$ ]]; then
  echo "Error: RELEASE_VERSION must be in X.Y format, but was $RELEASE_VERSION"
  exit 1
fi

set -x

docker compose build --pull
git tag -s -m "Release $RELEASE_VERSION" "v$RELEASE_VERSION"
docker tag luontola/gcp-dynamic-dns "luontola/gcp-dynamic-dns:$RELEASE_VERSION"
docker push luontola/gcp-dynamic-dns
docker push "luontola/gcp-dynamic-dns:$RELEASE_VERSION"
git push origin master "v$RELEASE_VERSION"
