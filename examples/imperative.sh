#!/bin/bash

# It's not currently sensible to actually run this script.
exit 1

set -e

function exop() {
  echo exop "$@"
  go run ./cmd/exop "$@"
}

function logp() {
  echo logp "$@"
  go run ./cmd/logp "$@"
}

exop /describe-components

exop post /delete

exop /create-component \
  name=tick \
  type=process \
  'spec={"program": "./bin/tick"}'

exop /create-component \
  name=echo \
  type=process \
  'spec={
    "program": "socat",
    "arguments": ["TCP4-LISTEN:2000,fork", "EXEC:cat"]
   }'

logp /

exo new process tick ./bin/tick

exo new container echo -p 2222:80 ealen/echo-server:0.5.1
