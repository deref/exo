#!/bin/env bash
set -euxo pipefail

source "$( dirname "${BASH_SOURCE[0]}" )/.include"

mkdir -p completions
for sh in bash zsh fish; do
	go run main.go completion generate "$sh" >"completions/exo.$sh"
done
