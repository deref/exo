#!/usr/bin/env bash

## Updates the version for release.

## TODO: Disallow when there are staged changes

source source "$( dirname "${BASH_SOURCE[0]}" )/.include"

currentversion="$(cat "${ROOT_DIR}/VERSION")"
version="$(date -u +'%Y.%m.%d')"
if [[ "$currentversion" == "${version}"* ]]; then
    lastbuild="${currentversion#*_}"
    if [[ "$lastbuild" == "$currentversion" ]]; then
        version="${version}_1"
    else
        version="${version}_$((lastbuild+1))"
    fi
fi

echo -n "$version" > "${ROOT_DIR}/VERSION"

