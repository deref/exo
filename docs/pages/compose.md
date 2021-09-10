# Compose

[Compose files][1] are a standard for the definition of multi-container
platform-agnostic applications. The standard originated with [Docker
Compose][2].

The simplest and most common Docker Compose workflow is the `up`
command:

```bash
docker-compose up
```

The analogous command with exo is `run`:

```bash
exo run
```

Both with Docker Compose and exo, this will start all the services defined in
your compose file and then tail their logs until interrupted. If you prefer
"detached mode" (via the `--detach` or `-d` flags), the analogous exo command
is `apply`:

```bash
exo apply
```

Most Docker Compose commands have an analogous exo command, often with the same
name! Peruse `exo help` for a list.

Please note that Compose compatibility in exo is currently experimental, and
not all features are supported yet. However, most simple Compose files should
work as they do with Docker Compose. If you encounter something that doesn't
work, please [let us know](https://github.com/deref/exo/issues)! We're working
hard to acheive feature parity.

## References

- [Compose Specification][1]
- [Docker Compose][2]

[1]: https://github.com/compose-spec/compose-spec
[2]: https://docs.docker.com/compose/
