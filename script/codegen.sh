#!/bin/bash

set -e

(
  cd gui
  npm run codegen
)

go run ./cmd/codegen
