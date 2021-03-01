#!/usr/bin/env bash

set -e

TARGET=$1

PANDO_VERSION=$(git describe --tags --abbrev=0)
PANDO_COMMIT=$(git rev-parse --short HEAD)

export PANDO_VERSION
export PANDO_COMMIT

ENVS_GEN="$TARGET/envs_gen.go"
trap 'rm -f $ENVS_GEN' EXIT
if ! type "envs" > /dev/null 2>/dev/null; then
  env GO111MODULE=off go get -u github.com/yiplee/envs
fi
envs --prefix PANDO -o "$ENVS_GEN"

CGO_ENABLED=0 go build -o builds/ "$TARGET"
