# Storage

The stateful components of exo typically store their state in a local document database. Deref has developed its own embedded database on top of [BadgerDB](https://github.com/dgraph-io/badger/), which provides a table-based API with many of the niceties of a SQL database.

The philosophy of our database is similar to that of FoundationDB in that there are multiple "layers" to the database, each of which exposes its own API, and client code can interact with whichever layer makes the most sense. The lowest level API simply gives K/V access, and the highest-level API provides a SQL-like DBMS.

## Tables

All data is stored in tables. A table is a collection of rows that are stored ordered by a primary key. Each table is assigned a 32-bit `ObjectID`, which functions as a key prefix for all items stored in that table.

The storage system also keeps track of several system tables that manage tables, indexes, and schemas. These tables have well-known identifiers that can be used by the lower levels of the storage engine directly. For example, to look up data in a user table called "log-events", the storage engine first looks up the "tables" collection using a well-known identifier.
