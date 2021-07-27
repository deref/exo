package storage

type ReadTransaction interface {
	// End closes a read transaction.
	End() error
}

type WriteTransaction interface {
	// Commit ensures that any changes made in the transaction are persisted and
	// are made visible to any future transactions.
	Commit() error

	// Rollback undoes the result of a transaction.
	// It must be a no-op to call Rollback() after committing a transaction.
	Rollback() error
}
