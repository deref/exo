package storage

import "fmt"

type Database struct {
	store KVEngine
}

func NewDatabase(store KVEngine) (*Database, error) {
	// TODO: Initialize internal schema.
	if ok, err := isBootstrapped(store); err != nil {
		return nil, fmt.Errorf("checking if database is bootstrapped: %w", err)
	} else if !ok {
		if err := bootstrapStore(store); err != nil {
			return nil, fmt.Errorf("bootstrapping database: %w", err)
		}
	}

	return &Database{
		store: store,
	}, nil
}

func (db Database) InstallTable(tbl *table) error {
	panic("TODO: Install table")
}
