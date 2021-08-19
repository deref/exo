# Changelog

All notable changes to this project will be documented in this file.

This project adhears to [CalVer](./doc/versioning.md).

## Unreleased

Docker Compose compatibility BETA.

### Added

- Dark Mode


## 2021-08-17

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


## 2021-07-29

Initial release.

### Added

- Procfile runner. Acts as a drop-in replacement for `foreman` with terminal
  and browser-based log viewers.
