# Testing

We don't have a lot of test automation yet, so here's a manual test plan.

# Features

## Installation

- `exo upgrade`
- `exo uninstall` and then [Fresh Installation](./install.md)

## Run

From two separate example directories:

- `exo run Procfile`
- `exo run compose.yaml`

This will create two separate workspaces.

Exercises: `exo apply`, `exo logs`, and import for supported manifest formats.

## CLI

- `exo start`, `exo stop`, `exo restart`
- `exo ps`, `exo ls`

## GUI

- Launch with `exo gui` in a workspace directory.
- Start/stop processes
- Toggle log viewers
- Switching between workspaces
- Create new process from GUI (`tick` is a good candidate)

# Performance and Robustness

## Logs

- Load test log throughput. Try [../examples/voluminous-logs] with both `exo
  logs` and log viewer in the GUI.

