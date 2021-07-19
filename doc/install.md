# Installation

## Quick Install

```bash
curl -sL https://exo.deref.io/install | sh
```

## Manual Installation

Grab the latest release from Github: https://github.com/deref/exo/releases

Put the binary `~/.exo/bin` and add that directory to your `PATH`.

## Uninstall

Kill and supervised processes:

```bash
exo destroy
```

Delete the exo home directory:

```bash
rm -rf ~/.exo
```

Remove `~/.exo/bin` from your path.
