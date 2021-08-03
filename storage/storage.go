package storage

type KVEngine interface {
	Scanner

	Set(tx WriteTransaction, key, val []byte) error
	Get(tx ReadTransaction, key []byte) ([]byte, error)

	ReadTransaction() ReadTransaction
	WriteTransaction() WriteTransaction
}

func SetAtomic(kv KVEngine, key, val []byte) error {
	return Transact(kv, func(txn WriteTransaction) error {
		return kv.Set(txn, key, val)
	})
}

func GetAtomic(kv KVEngine, key []byte) ([]byte, error) {
	var val []byte
	err := View(kv, func(txn ReadTransaction) (getErr error) {
		val, getErr = kv.Get(txn, key)
		return
	})
	return val, err
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
