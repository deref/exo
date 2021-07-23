# exo: a process manager & log viewer for dev

**exo-** _prefix_ – external; from outside.

![The Exo GUI](https://github.com/deref/exo/blob/main/doc/screenshot-light.png?raw=true)

## Status

_Alpha!_

## Getting Started

Install exo:

```bash
curl -sL https://exo.deref.io/install | sh
```

If you prefer manual installation, see [./doc/install.md](./doc/install.md) for
details, including uninstall instructions.

Navigate to your code directory and then launch the exo gui:

```bash
exo gui
```

For more features, consult the builtin help:

```bash
exo help
```

---

## Telemetry

**exo** collects limited and anonymous telemetry data by default. This behavior
can be disabled by adding the following setting to your exo config (located at
`~/.exo/config.toml` by default):

```bash
[telemetry]
disable: true
```

