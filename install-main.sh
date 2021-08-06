#!/usr/bin/env bash

set -e

function build_ui {
  (
    cd gui
    npm i
    npm run build
  )
}

BIN_DIR="${HOME}/.exo/bin"
DEV_BIN="${BIN_DIR}/exo_dev"
EXO_LINK="${BIN_DIR}/exo"
build_ui
go build -tags bundle -o "$DEV_BIN" ./cmd/exo
ln -sf "$DEV_BIN" "$EXO_LINK"
