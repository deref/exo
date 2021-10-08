#!/bin/bash

set -ex

which exo && exo exit || true
go build -o ./bin/exo
mkdir -p ~/.exo/bin
cp ./bin/exo ~/.exo/bin/exo
