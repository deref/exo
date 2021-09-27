#!/bin/bash

set -ex

node --loader ts-node/esm ./script/gen-icon-list.ts
node --loader ts-node/esm ./script/gen-theme.ts
