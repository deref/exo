# Log Deamon

This is a simple log collection service for fifo pipe sources.

Collector processes are pseudo-stateless. That is, all state is externalized to
the file system with the exception of collection streaming workers. These
workers are designed to be safe in the event that more than one worker is
running against the same log source (although this is, as of writing,
unverified).

# Resources

## Logs

CRUD operations: `/add-log`, `/remove-log`, and `/describe-logs`.

## Events

Each log is made up of events, which combine a log line message with a globally unique ID that is guaranteed to increase monotonically.

IDs are 26-character base32-encoded (lowercase) strings.

The `/get-events` operation provides paginated queries of the union of
several log streams.

### Collecting

In normal operation, the `logd` deamon collects logs by reading from the
source fifos, decorating each message metadata, and then
writing the event to BadgerDB.

During development, the `/collect` operation can be invoked directly.

# Operation

Use the `logd` command.

## Log Rotation

*NOT YET IMPLEMENTED*

# Development

See `logp` command.
