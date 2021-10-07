#!/usr/bin/env bash

##
## Builds the current development server and GUI in release mode and
## installs it to ~/.exo/bin/exo.
##

set -e

while [[ $# -gt 0 ]]; do
  flag=$1
  shift
  case $flag in
    --skip-build-gui)
      skip_build_gui=1
      ;;
  esac
done

which exo && exo exit || true

function build_gui {
  (
    cd gui
    npm i
    npm run build
  )
}

BIN_DIR="${HOME}/.exo/bin"
DEV_BIN="${BIN_DIR}/exo_dev"
EXO_LINK="${BIN_DIR}/exo"
if [[ ! $skip_build_gui ]]; then
  build_gui
fi
go build -tags bundle -o "$DEV_BIN"
ln -sf "$DEV_BIN" "$EXO_LINK"

