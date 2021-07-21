# Release Process

exo uses GitHub actions to publish a new release of exo any time the `core/VERSION` file changes. exo follows a (CalVer)[https://calver.org/] versioning scheme in which each release version adheres to the following scheme: `YYYY.MM.DD_MICRO` where the `_MICRO` suffix is only present if there have been multiple releases on a single day.

To trigger a new build, increment the version with: `./scripts/increment-version.sh`.
