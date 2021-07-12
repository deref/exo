# Log Collector

This is a simple log collector service for fifo pipe sources.

Collector processes are pseudo-stateless. That is, all state is externalized to
the file system with the exception of collection streaming workers. These
workers are designed to be safe in the event that more than one worker is
running against the same log source (although this is, as of writing,
unverified).

# Resources

## Logs

CRUD operations: `/add-log`, `/remove-log`, and `/describe-logs`.

## Events

Each log is made up of events, which combine a log line message with a per-log
sequence id (abbreviated "sid") and a timestamp.

The `/get-events` operation provides paginated queries of the union of
several log streams. The log storage format is indexed by sid.

TODO: Scan with before/after parameters for both sid and timestamp ranges.

### Collecting

In normal operation, the `logd` deamon collects logs by reading from the
source fifos, decorating each message with a sid and timestamp, and then
writing the event to a log file.

During development, the `/collect` operation can be invoked directly.

# Operation

Use the `logd` command.

## Log Rotation

*NOT YET IMPLEMENTED*

# Development

See `logp` command.
