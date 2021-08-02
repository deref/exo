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

func (db *Database) Table(name string) (tbl *table, err error) {
	tbl = &table{
		db: db,
	}

	tt := tableTable(db)
	var tblTup *Tuple
	if tblTup, err = SelectOne(tt, ColumnByNameEquals(tt.schema, "table_name", name)); err != nil {
		return nil, err
	} else if tblTup == nil {
		return
	}

	oid, _ := tblTup.GetInt32(0) // TODO: Use a GetXXXByName method.
	tbl.oid = OID(oid)

	var schemaOID int32
	schemaOID, _ = tblTup.GetInt32(1)

	tbl.name, _ = tblTup.GetUnicode(2)

	st := schemaTable(db)
	var schemaTup *Tuple
	if schemaTup, err = SelectOne(st, ColumnByNameEquals(st.schema, "schema_oid", schemaOID)); err != nil {
		return nil, err
	} else if schemaTup == nil {
		return
	}
	schemaData, err := schemaTup.GetBytes(1)
	if err != nil {
		return nil, err
	}
	tbl.schema = MustDeserializeSchema(schemaData)
	tbl.serde = NewSchematizedRowSerde(tbl.schema)

	return tbl, nil
}
