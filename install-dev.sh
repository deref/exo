#!/bin/bash

set -ex

exo exit || true
go build -o ./bin/ ./cmd/exo
mkdir -p ~/.exo/bin
cp ./bin/exo ~/.exo/bin/exo
