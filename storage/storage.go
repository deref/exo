package storage

type KVEngine interface {
	Scanner

	Set(key, val []byte) error
	Get(key []byte) ([]byte, error)

	ReadTransaction() ReadTransaction
	WriteTransaction() WriteTransaction
}

type TransactFunc = func(txn WriteTransaction) error

func Transact(kv KVEngine, fn TransactFunc) error {
	txn := kv.WriteTransaction()
	defer txn.Rollback()
	if err := fn(txn); err != nil {
		return err
	}

	return txn.Commit()
}

type ViewFunc = func(txn ReadTransaction) error

func View(kv KVEngine, fn ViewFunc) error {
	txn := kv.ReadTransaction()
	defer txn.End()
	return fn(txn)
}
