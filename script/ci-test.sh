#!/bin/bash

set -euxo pipefail

(cd gui && npm ci)

./script/codegen.sh

if [[ "$(git status --porcelain)" ]]; then
  echo "Regenerate script modified files. Please run ./script/codegen.sh"
  echo "These are the changes:"
  git diff
  exit 1
fi

go test ./...

(cd gui && npm run check)
