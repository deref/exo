#!/usr/bin/env bash

source "$( dirname "${BASH_SOURCE[0]}" )/.include"

export EXO_HOME="$EXO_DEV_HOME"

go run ./cmd/watch --dir "${ROOT_DIR}" -r yes --ignore '.git,var,gui,.dev,examples' -- \
  go run "${ROOT_DIR}" server --force-std-log=1
