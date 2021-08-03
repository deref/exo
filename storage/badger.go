package storage

import (
	"fmt"

	"github.com/dgraph-io/badger/v3"
)

func NewMemoryKVEngine() *BadgerKVEngine {
	opts := badger.DefaultOptions("").WithInMemory(true)
	db, err := badger.Open(opts)
	if err != nil {
		panic(fmt.Errorf("Opening in-memory database should never fail: %w", err))
	}

	return &BadgerKVEngine{
		db: db,
	}
}

type BadgerKVEngine struct {
	db *badger.DB
}

var _ KVEngine = (*BadgerKVEngine)(nil)

func (kv *BadgerKVEngine) Get(txn ReadTransaction, key []byte) ([]byte, error) {
	tx := txn.(*badgerReadTransaction).inner
	item, err := tx.Get(key)
	switch err {
	case badger.ErrKeyNotFound:
		return nil, nil
	case nil:
		// Do nothing.
	default:
		return nil, err
	}

	return item.ValueCopy(nil)
}

func (kv *BadgerKVEngine) Set(txn WriteTransaction, key, val []byte) error {
	tx := txn.(*badgerWriteTransaction).inner
	return tx.Set(key, val)
}

func (kv *BadgerKVEngine) ReadTransaction() ReadTransaction {
	return &badgerReadTransaction{
		inner: kv.db.NewTransaction(false),
	}
}

type badgerReadTransaction struct {
	inner *badger.Txn
}

func (txn *badgerReadTransaction) End() error {
	txn.inner.Discard()
	return nil
}

func (kv *BadgerKVEngine) WriteTransaction() WriteTransaction {
	return &badgerWriteTransaction{
		inner: kv.db.NewTransaction(true),
	}
}

type badgerWriteTransaction struct {
	inner *badger.Txn
}

func (txn *badgerWriteTransaction) Rollback() error {
	txn.inner.Discard()
	return nil
}

func (txn *badgerWriteTransaction) Commit() error {
	return txn.inner.Commit()
}

func (txn *badgerWriteTransaction) Downgrade() ReadTransaction {
	return &badgerReadTransaction{
		inner: txn.inner,
	}
}

func (kv *BadgerKVEngine) Scan(txn ReadTransaction, args ScanArgs) ScanIter {
	opts := badger.DefaultIteratorOptions
	opts.Prefix = args.Prefix
	if args.Direction == ScanDirectionDESC {
		opts.Reverse = true
	}
	if args.KeyOnly {
		opts.PrefetchValues = false
	}

	it := txn.(*badgerReadTransaction).inner.NewIterator(opts)
	it.Seek(opts.Prefix)

	return &badgerScanIter{
		inner:   it,
		prefix:  opts.Prefix,
		keyOnly: args.KeyOnly,
	}
}

type badgerScanIter struct {
	inner   *badger.Iterator
	prefix  []byte
	keyOnly bool

	next *badger.Item
	err  error
}

func (it *badgerScanIter) Next() bool {
	if !it.inner.ValidForPrefix(it.prefix) {
		return false
	}
	it.next = it.inner.Item()
	it.inner.Next()
	return true
}

func (it *badgerScanIter) Item() *ScanEntry {
	item := it.next
	key := item.Key()
	entry := &ScanEntry{
		Key: key,
	}
	if it.keyOnly {
		return entry
	}

	entry.Value, it.err = item.ValueCopy(nil)

	return entry
}

func (it *badgerScanIter) Err() error {
	return it.err
}

func (it *badgerScanIter) Close() {
	it.inner.Close()
}
