#!/usr/bin/env bash
set -euo pipefail
NAME="${1:?name required}"
PLATFORM="${2:-linux/amd64}"
docker build --platform "${PLATFORM}" -f benzhi.Dockerfile -t "benzhi/${NAME}:latest" .
