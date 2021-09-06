# Uninstallation

## Automated Uninstall

```bash
exo uninstall
```

## Manual Uninstall

If you installed shell completions, reference `exo completion --help` for
information on relevant install paths for your preferred shell.

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
