# Changelog

All notable changes to this project will be documented in this file.

This project adhears to [CalVer](./doc/versioning.md).

## [Unreleased]

### Added

- Linkify URLs in log messages.
- Support colors and other text attributes in log messages.
- Docker Compose compatibility.
- Job system for reporting parallel task progress.

## Changed

- Overhauled look and feel.

### Fixed

- Shell syntax handling in Procfiles.
- GUI fails gracefully when backend is down.
- Terminate whole process groups more reliably.


## 2021-07-29

Initial release.

### Added

- Procfile runner. Acts as a drop-in replacement for `foreman` with terminal
  and browser-based log viewers.
