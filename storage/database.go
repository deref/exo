package storage

import (
	"fmt"
	"sync"
)

type Database struct {
	store KVEngine

	mu      sync.Mutex
	nextOID OID
	// cachedTables map[string] *table // TODO: Cache tables by name for loading w/o going to schema table.
}

func NewDatabase(store KVEngine) (*Database, error) {
	db := &Database{
		store: store,
	}

	if ok, err := db.isBootstrapped(); err != nil {
		return nil, fmt.Errorf("checking if database is bootstrapped: %w", err)
	} else if !ok {
		if err := db.bootstrap(); err != nil {
			return nil, fmt.Errorf("bootstrapping database: %w", err)
		}
	}

	return db, nil
}

func (db *Database) getNextOID() OID {
	db.mu.Lock()
	defer db.mu.Unlock()
	oid := db.nextOID
	db.nextOID++
	return oid
}

func (db *Database) ReadTransaction() ReadTransaction {
	return db.store.ReadTransaction()
}

func (db *Database) WriteTransaction() WriteTransaction {
	return db.store.WriteTransaction()
}

func (db *Database) Tables() *table {
	return tableTable(db)
}

func (db *Database) Scemas() *table {
	return schemaTable(db)
}

func (db *Database) Table(txn ReadTransaction, name string) (*table, error) {
	tbl := &table{
		db: db,
	}

	tt := tableTable(db)
	var tblTup *Tuple
	tblTup, err := SelectOne(txn, tt, ColumnByNameEquals(tt.schema, "table_name", name))
	if err != nil {
		return nil, err
	} else if tblTup == nil {
		return nil, nil
	}

	oid, _ := tblTup.GetInt32(0) // TODO: Use a GetXXXByName method.
	tbl.oid = OID(oid)

	var schemaOID int32
	schemaOID, _ = tblTup.GetInt32(1)

	tbl.name, _ = tblTup.GetUnicode(2)

	st := schemaTable(db)
	var schemaTup *Tuple
	if schemaTup, err = SelectOne(txn, st, ColumnByNameEquals(st.schema, "schema_oid", schemaOID)); err != nil {
		return nil, err
	} else if schemaTup == nil {
		return nil, nil
	}
	schemaData, err := schemaTup.GetBytes(1)
	if err != nil {
		return nil, err
	}
	tbl.schema = MustDeserializeSchema(schemaData)
	tbl.serde = NewSchematizedRowSerde(tbl.schema)

	return tbl, err
}
