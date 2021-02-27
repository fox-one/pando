#!/usr/bin/env bash

set -e

TARGET=$1

PANDO_VERSION=$(git describe --tags)
PANDO_COMMIT=$(git rev-parse --short HEAD)

export PANDO_VERSION
export PANDO_COMMIT

ENVS_GEN="$TARGET/envs_gen.go"
if ! type "envs" > /dev/null; then
  go install github.com/yiplee/envs@latest
fi
envs --prefix PANDO -o "$ENVS_GEN"

CGO_ENABLED=0 go build -o builds/ "$TARGET"

# clean envs
rm "$ENVS_GEN"
