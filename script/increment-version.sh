#!/usr/bin/env bash

## Updates the version for release.

## TODO: Disallow when there are staged changes

ROOTDIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )/.." &> /dev/null && pwd )"

currentversion="$(cat "${ROOTDIR}/VERSION")"
version="$(date -u +'%Y.%m.%d')"
if [[ "$currentversion" == "${version}"* ]]; then
    lastbuild="${currentversion#*_}"
    if [[ "$lastbuild" == "$currentversion" ]]; then
        version="${version}_1"
    else
        version="${version}_$((lastbuild+1))"
    fi
fi

echo -n "$version" > "${ROOTDIR}/VERSION"

