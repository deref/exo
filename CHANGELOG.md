# Changelog

All notable changes to this project will be documented in this file.

This project adhears to [CalVer](./doc/versioning.md).

## Unreleased

### Added

- Project templates


## 2021.10.06

### Added

- Ability to create new projects in the GUI.
- GUI for deleting components.
- `exo edit` command for editing component specs.
- `exo rename` command to rename components.
- Commands to clear log scrollback in CLI and GUI.
- Workspace log stream for system events in log viewer.
- The `exo gui` command learned the `--print` flag.
- `exo exec` command runs one-off processes in the workspace Environment.
- Workspace information view in GUI.

### Fixed

- Numerous issues applying modified manifests.
- Various error reporting improvements.
- Migrated log storage from BadgerDB to Sqlite.
- Secured CORS policy and daemon request authorization.
- Secured iframe embedding against click-jacking attacks.
- Secured against DNS rebinding attacks.
- Workspaces with the same last directory path component are no longer
  ambiguous.
- Enhanced appearance of process list.
- Unintended disabling of colored logging with popular JavaScript and other
  libraries.
- Unbounded growth of log storage.


## 2021.09.14

### Added

- Status and progress of docker image pulls
- [#318](https://github.com/deref/exo/pull/318) Aligned job messages
- [#321](https://github.com/deref/exo/pull/321) Add command to edit raw specification

### Fixed

- [#316](https://github.com/deref/exo/pull/316) Shell form docker commands when image pull policy set
- [#323](https://github.com/deref/exo/pull/323) Duplicated new lines
- [#327](https://github.com/deref/exo/issues/327) Characters omitted in log output

## 2021.09.08

## Fixed

- Various docker-compose parsing issues.


## 2021.09.07

### Added

- [#313](https://github.com/deref/exo/pull/213) Shell completions (optional install).
- [#284](https://github.com/deref/exo/pull/284) `exo env` command.
- Additional docker-compose compatibility: anchors & aliases, x- fields, cpu
  contraints, arbitrary CMD syntax, and more.
- [#280](https://github.com/deref/exo/pull/280) `exo state ...` commands for performing state repair.
- [#188](https://github.com/deref/exo/issues/188) `exo kill -s signal ...` command.

### Fixed

- [#273](https://github.com/deref/exo/pull/273) Truncated output after processes terminate.
- [e062d58](https://github.com/deref/exo/commit/e062d589fec56fcbefc777444eb6d1ac4ddf0d7d) Fix `run` exiting prematurely on some Mac systems.


## 2021.08.31

### Fixed

- Docker compose volume name prefixing.

## 2021.08.28

### Fixed

- Directory creation race on first run.


## 2021.08.27

Docker Compose compatibility BETA.

### Added

- Dark Mode
- Component dependencies
- Support for `.env` files

### Fixed

- Many Docker compose compatibility issues.


## 2021.08.17

Docker Compose compatibility ALPHA.

### Added

- Linkify URLs in log messages.
- Support colors and other text attributes in log messages.
- Job system for reporting parallel task progress.
- exo²: Use the current version of exo to develop the next version of exo.
- Log filtering in GUI.
- Process status details screens.
- Overhauled look and feel.

### Fixed

- Shell syntax handling in Procfiles.
- GUI fails gracefully when backend is down.
- Terminate whole process groups more reliably.
- Log viewer line wrapping.


## 2021.07.29

Initial release.

### Added

- Procfile runner. Acts as a drop-in replacement for `foreman` with terminal
  and browser-based log viewers.
