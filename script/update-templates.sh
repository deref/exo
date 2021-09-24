#!/usr/bin/env bash
set -euo pipefail

# Call this script to update the set of templates stored in the S3 starter
# templates bucket. This should be called every time the underlying templates
# repository is updated. In order to run successfully, one will need the AWS
# CLI to be set up with credentials that grant write permission to the
# exo-starter-templates S3 bucket.

source "$( dirname "${BASH_SOURCE[0]}" )/.include"

cd "$ROOT_DIR"
aws s3 sync "$(go run ./cmd/gentemplates/main.go)" s3://exo-starter-templates
