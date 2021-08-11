# Local Development

exo is a client/server application whose client is written using TypeScript and Svelte and whose server is written in Go. In order to develop exo, you must have the standard development tools for these platforms installed:

- [Go >= 1.16](https://golang.org/doc/install).
- [Node.js >= 14.x](https://nodejs.org/en/download/) with NPM >= 6.x (typically included with Node.js).
- [exo >= 2021.08.04](https://exo.deref.io). We use a released version of exo to develop prerelease version.

The exo repository contains a Procfile that will start the exo server and the exo GUI in development mode. To create a workspace for developing exo, please run the following:

```bash
cd path/to/exo
./script/deps.sh
exo run
```

If all goes well, you should be able to manage your server and gui processes at [http://localhost:4000](http://localhost:4000). Please note that this is the _installed_ exo gui that you are viewing, not the development instance. The development instance runs on port `4001` and can be accessed at [http://localhost:4001](http://localhost:4001). The development mode GUI has a "DEV" indicator in the footer so that you can tell at a glance which instance you are using.

The `exo` CLI runs against the installed instance by default, but you can change to the development instance by adding the following to your exo config file (located at `~/.exo/config.toml`):

```
[client]
url = "http://localhost:4001"
```

Now all `exo` commands will run against the development instance. To run against the installed instance again, remove or comment out these lines from your `config.toml`.

## envrc

The project root directory contains an `.envrc` script. You can use
[https://direnv.net/](direnv) to automatically source this script when
working on the project.

## dexo

To test CLI changes against the auto-restarting watched development server, use
`./script/dexo`. If you're using direnv (see above), then `dexo` will be on
your `PATH`.

### Troubleshooting

#### The GUI won't start. I see an "Error: esbuild: Failed to install correctly" error in the log.

Manually install esbuild with the following command:

```bash
node node_modules/esbuild/install.js
```

Then restart the GUI.
