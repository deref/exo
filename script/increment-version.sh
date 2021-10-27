#!/usr/bin/env bash

## Updates the version for release.

## TODO: Disallow when there are staged changes

source "$( dirname "${BASH_SOURCE[0]}" )/.include"

currentversion="$(cat "${ROOT_DIR}/VERSION")"
version="$(date -u +'%Y.%m.%d')"
if [[ "$currentversion" == "${version}"* ]]; then
    lastbuild="${currentversion#*-}"
    if [[ "$lastbuild" == "$currentversion" ]]; then
        version="${version}-1"
    else
        version="${version}-$((lastbuild+1))"
    fi
fi

echo -n "$version" > "${ROOT_DIR}/VERSION"

