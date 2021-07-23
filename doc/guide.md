# Guide

## Concepts

**Workspaces** - A mapping of filesystem paths. Most projects have one
workspace rooted at the same directory as their checked out code. This is how
`exo` knows what project you're working on based on your current working
directory.

If run in an unmapped directory, `exo gui` will offer to create a workspace for
you. You can determine the current workspace with `exo workspace`, initialize a
new one with `exo workspace init` or delete the current workspace with
`exo workspace destroy`.

**Components** - An abstract definition of resources managed by exo. Presently,
the only supported type of components are _processes_. Each component has a
unique name within a workspace. Components are manipulated by editing applying
manifests (see below), or with CRUD operations such as `exo ls`, `exo new`, and
`exo rm`.

**Manifests** - A file describing all of the components in a project.
Presently, [procfiles](./procfiles.md) are the only supported type of
manifest. Use `exo apply ./path/to/Procfile` whenever your Procfile changes
to make your workspace match. Processes will be added or removed accordingly.

**Processes** - A running program. Presently, only host-machine processes are
supported. Docker containers will be supported in the future.

Assuming you have a process named `myapp`, here are some useful management
commands:

```bash
# List processes.
exo ps

# Control process state.
exo start myapp
exo restart myapp
exo stop myapp
```

**Daemon** - A background service that manages components and supervises
processes. Most commands start this service automatically. You can start it
explicitly with `exo daemon` and terminate it with `exo exit`.

## Workflow

For most standard procfile setups, `exo gui` will do the right thing on first
use automatically. You can manage processes and view logs in your browser.

If you've got multiple manifests, such as different Procfiles for both dev and
test. Or if you generally prefer command line interfaces, a typical workflow
looks something like this:

```bash
# Initialize a new workspace in the current directory.
exo workspace init

# Apply a manifest to start it's processes.
exo apply ./Procfile.dev

# Tail logs in your terminal.
exo logs

# Or only specific processes.
exo logs api worker

# Manipulate individual processes.
exo stop worker
exo restart api

# Switch to a different configuration by applying a different manifest.
exo apply ./Procfile.test

# Shutdown everything and cleanup state when you're done.
exo workspace destroy

# If you're very, very done and don't want exo running anymore.
exo exit
```
