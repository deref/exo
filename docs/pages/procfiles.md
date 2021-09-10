# Procfiles

Procfiles are a simple manifest format that describes processes. Here's an
example:

```Procfile
web: node run dev
api: go run ./server
```

Each line contains a process name and a command to execute separated by a colon
and some whitespace.

When a process is started, a unique `PORT` environment variable is supplied to
each.

Most procfile runners, exo included, have a single command to start all
processes in the procfile, then tail their logs until interrupted:

```bash
exo run Procfile
```

## References

- [The New Heroku](https://blog.heroku.com/the_new_heroku_1_process_model_procfile) - Introduction of Procfiles from 2011.
- [Foreman](https://github.com/ddollar/foreman) - The original Procfile runner.
