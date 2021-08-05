#!/usr/bin/env bash

source "$( dirname "${BASH_SOURCE[0]}" )/.include"

env EXO_HOME="$EXO_DEV_HOME" \
    go run ./cmd/watch --dir "${ROOT_DIR}" -r yes --ignore '.git,var,gui' -- \
        go run "${ROOT_DIR}/cmd/exo" server --force-std-log=1
