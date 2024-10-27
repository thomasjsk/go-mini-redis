#!/bin/sh
set -e

(
  cd "$(dirname "$0")"
  go build -o /tmp/go-mini-redis app/*.go
)

exec /tmp/go-mini-redis "$@"
