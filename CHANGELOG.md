# Changelog

All notable changes to this project will be documented in this file.

This project adhears to [CalVer](./doc/versioning.md).

## Unreleased

### Added

- Project templates

## 2021.10.29

### Fixed

- [#77](https://github.com/deref/exo/pull/77) Handling of shell-syntax in
  Procfiles.
- [#459](https://github.com/deref/exo/pull/459) Compose file parsing of !!int
  tag in string positions.

## 2021.10.28

### Added

- [#393](https://github.com/deref/exo/pull/393), [#456](https://github.com/deref/exo/pull/456) Initial integration with Deref's secrets service.

## 2021.10.27

### Added

- [#420](https://github.com/deref/exo/pull/420) Support for variable interpolation in compose files.
- [#418](https://github.com/deref/exo/pull/418), [#444](https://github.com/deref/exo/pull/444), [#445](https://github.com/deref/exo/pull/445) Initial support for `exo.hcl` files.

### Changed

- [#429](https://github.com/deref/exo/pull/393) Now licensed as Apache v2.

### Fixed

- [#434](https://github.com/deref/exo/pull/393) Colors not synchronized between components list and log viewer.

## 2021.10.12

- [#407](https://github.com/deref/exo/pull/407) `exo init` command.
- [#406](https://github.com/deref/exo/pull/406) Alpha/undocumented release of `exo.hcl` manifests.

## 2021.10.08_1

### Added

- [#401](https://github.com/deref/exo/pull/401) Ability to edit component specs from the GUI.

### Fixed

- Race condition in token handling.
- Context menu behaviors.

## 2021.10.08

### Added

- [#337](https://github.com/deref/exo/pull/337) Ability to create new projects in the GUI.
- [#392](https://github.com/deref/exo/pull/392) GUI for deleting components.
- [#321](https://github.com/deref/exo/pull/321) `exo edit` command for editing component specs.
- [#400](https://github.com/deref/exo/pull/400) `exo rename` command to rename components.
- [#347](https://github.com/deref/exo/issues/347) Commands to clear log scrollback in CLI and GUI.
- [#98](https://github.com/deref/exo/issues/98) Workspace log stream for system events in log viewer.
- [#382](https://github.com/deref/exo/pull/382) The `exo gui` command learned the `--print` flag.
- [#371](https://github.com/deref/exo/pull/371) `exo exec` command runs one-off
  processes in the workspace Environment.
- [#369](https://github.com/deref/exo/pull/369) Workspace information view in GUI.

### Fixed

- Numerous issues applying modified manifests.
- Various error reporting improvements.
- [#358](https://github.com/deref/exo/pull/358) Migrated log storage from
  BadgerDB to Sqlite.
- [#360](https://github.com/deref/exo/issues/360) Secured CORS policy and
  daemon request authorization.
- [#377](https://github.com/deref/exo/pull/377) Secured iframe embedding
  against click-jacking attacks.
- [#361](https://github.com/deref/exo/pull/361) Secured against DNS rebinding attacks.
- [#339](https://github.com/deref/exo/pull/339) Workspaces with the same last
  directory path component are no longer ambiguous.
- [#335](https://github.com/deref/exo/pull/335) Enhanced appearance of process list.
- [#365](https://github.com/deref/exo/pull/365) Unintended disabling of colored
  logging with popular JavaScript and other libraries.
- [#336](https://github.com/deref/exo/issues/336) Unbounded growth of log storage.

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
