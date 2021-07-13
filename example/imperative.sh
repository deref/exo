#!/bin/bash

set -e

function exop() {
  echo exop "$@"
  go run ./cmd/exop "$@"
}

exop /describe-components

exop /create-component \
  name=echo \
  type=process \
  'spec:={
    "command": "socat",
    "arguments": ["TCP4-LISTEN:2000,fork", "EXEC:cat"]
   }'
