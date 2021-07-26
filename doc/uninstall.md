# Uninstallation

## Automated Uninstall

```bash
exo uninstall
```

## Manual Uninstall

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

