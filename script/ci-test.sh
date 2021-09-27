#!/bin/bash

set -e

# TODO: Configure Node support in CI and run ./script/codegen.sh instead.
go run ./cmd/codegen

if [[ "$(git status --porcelain)" ]]; then
  echo "Regenerate script modified files. Please run ./script/codegen.sh"
  echo "These are the changes:"
  git diff
  exit 1
fi

go test ./...
