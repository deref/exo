# Installation

## Quick Install

```bash
curl -sL https://exo.deref.io/install | sh
```

## Manual Installation

Grab the latest release from Github: https://github.com/deref/exo/releases

Put the binary `~/.exo/bin` and add that directory to your `PATH`.

## Uninstall

Stop supervised processes and free any workspace resources:

```bash
cd /your/workspace
exo workspace destroy
```

Have multiple workspaces? Find them all:

```bash
exo workspace ls
```

Shutdown the exo daemon:

```bash
exo exit
```

Delete the exo home directory:

```bash
rm -rf ~/.exo
```

Remove `~/.exo/bin` from your path.
