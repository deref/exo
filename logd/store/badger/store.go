package badger

import (
	"context"

	"github.com/deref/exo/gensym"
	"github.com/deref/exo/logd/store"
	"github.com/dgraph-io/badger/v3"
)

type Store struct {
	db    *badger.DB
	idGen *gensym.ULIDGenerator
}

type Log struct {
	db    *badger.DB
	idGen *gensym.ULIDGenerator
	name  string
}

func Open(ctx context.Context, logsDir string) (*Store, error) {
	db, err := badger.Open(badger.DefaultOptions(logsDir))
	if err != nil {
		return nil, err
	}
	return &Store{
		db:    db,
		idGen: gensym.NewULIDGenerator(ctx),
	}, nil
}

func (sto *Store) Close() error {
	return sto.db.Close()
}

func (sto *Store) GetLog(name string) store.Log {
	return &Log{
		db:    sto.db,
		idGen: sto.idGen,
		name:  name,
	}
}
