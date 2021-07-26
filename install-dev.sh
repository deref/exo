#!/bin/bash

set -ex

go build -o bin ./cmd/exo
mkdir -p ~/.exo/bin
cp ./bin/exo ~/.exo/bin/exo
