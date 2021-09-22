#!/usr/bin/env bash
set -euo pipefail

source "$( dirname "${BASH_SOURCE[0]}" )/.include"

cd "$ROOT_DIR"
aws s3 sync "$(go run ./cmd/gentemplates/main.go)" s3://exo-starter-templates
