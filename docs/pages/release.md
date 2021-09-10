# Release Process

exo uses GitHub actions to publish a new release of exo any time the `VERSION` file changes. exo follows a [CalVer](https://calver.org/) versioning scheme in which each release version adheres to the following scheme: `YYYY.MM.DD_MICRO` where the `_MICRO` suffix is only present if there have been multiple releases on a single day.

To trigger a new build, increment the version with: `./scripts/increment-version.sh`.

## Details

### CI

When the contents of `VERSION` changes on the `main` branch, the [Create a release on VERSION update](https://github.com/deref/exo/blob/main/.github/workflows/perform-release.yaml) GitHub Action is triggered. This action performs the following steps:

1. Run all tests.
2. Create a tag that will be applied to the release in the form `v<VERSION>`, e.g. `v2021.09.08_1`.
3. Create a GitHub release from that tag.
4. Compile exo for all supported platforms and architectures, uploading the binary and SHA256 checksum for each.
5. Update the version stored in a [CloudFlare Workers KV bucket](https://dash.cloudflare.com/b5d1a3835f0dd218f633c7d612935c6b/workers/kv/namespaces) (this link only works for Deref CloudFlare admins).

### CloudFlare Worker

The CloudFlare K/V bucket is read by the [download-page worker](https://dash.cloudflare.com/b5d1a3835f0dd218f633c7d612935c6b/workers/view/download-page), whose code is in the [exo-install](https://github.com/deref/exo-install) repository. This worker serves 2 endpoints:

- `/version.txt`: returns the current version as a plain text document.
- `/install.sh`: the installer script, whish uses a template to populate the current version.

In order to decouple our code from CloudFlare, we proxy these two endpoints at `https://exo.deref.io/install` and `https://exo.deref.io/latest-version` respectively.


### Installer

The installer served by the CloudFlare worker via `https://exo.deref.io/install` is designed to be piped to Bash, and it performs the following tasks:

1. Detect the user's platorm and architecture.
2. Download the appropriate binary from the most recent GitHub release.
3. Compute a checksum of the binary and compare to the appropriate checksum from the GitHub release.
4. Create a symlink to the version of the binary downloaded to `~/.exo/bin/exo`

This install script can be run directly using the following command:

```
curl -sL https://exo.deref.io/install | sh
```

### Update Process

If exo is already running, the daemon periodically checks `https://exo.deref.io/latest-version` for the most recent version. If this version is later than the version compiled into the running process, it indicates that an update is available. The GUI periodically checks whether an update is available, and if it is, it displays a notice indicating the new version and gives the user the option to update. When the user updates, the daemon downloads and runs the latest installer. When the installation is complete, the daemon uses `execve` to replace itself with the newly-downloaded version. When the GUI detects that the daemon has restarted and is now up to date, it forces a reload of the page so that the GUI will also be updated.
