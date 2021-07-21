# Versioning

## Core Components

exo's core components use [CalVer](https://calver.org/).

Format is `YYYY.0M.0D_MICRO` with the following components:

- `YYYY` - four digit year
- `0M` - zero-padded, two-digit, month number
- `0D` - zero-padded, two-digit, day of month
- `_MICRO` - optional. single-digit, intra-day patch number.

The patch number is omitted for the first release on any given day.

## File Formats

Data formats are versioned using a SemVer-compatible scheme extended with
guidance for stable migration paths between major version.

See https://gist.github.com/brandonbloom/465625acaf0120354614e7fc0c117c62 for
details.
