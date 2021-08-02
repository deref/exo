package storage

import (
	"errors"
	"fmt"
)

type TableOptions struct {
	OID     *OID
	Name    string
	Schema  *Schema
	Indexes []IndexOptions
}

func (db *Database) NewTable(opts TableOptions) *table {
	t := &table{
		db:     db,
		name:   opts.Name,
		schema: opts.Schema,
		serde:  NewSchematizedRowSerde(opts.Schema),
		// TODO: Indexes.
	}

	if opts.OID != nil {
		t.oid = *opts.OID
	}

	return t
}

type table struct {
	db    *Database
	valid bool

	oid    OID
	name   string
	schema *Schema
	// XXX: A serde is not serializable, so the table definition that is saved to the db
	// should have a reference to a serde that the application is aware of.
	serde Serde
}

func (tbl *table) IsValid() bool {
	return tbl.oid != 0
}

func (tbl *table) IsSystem() bool {
	return tbl.oid < 0
}

func (tbl *table) Create() error {
	if tbl.db == nil {
		return errors.New("database not set")
	}
	if tbl.oid == 0 {
		tbl.oid = tbl.db.getNextOID()
	}

	schemaOID := tbl.db.getNextOID()
	if err := schemaTable(tbl.db).Insert(map[string]interface{}{
		"schema_oid":        int32(schemaOID),
		"serialized_schema": MustSerializeSchema(tbl.schema),
	}); err != nil {
		return fmt.Errorf("inserting schema: %w", err)
	}

	if err := tableTable(tbl.db).Insert(map[string]interface{}{
		"table_oid":  int32(tbl.oid),
		"schema_oid": int32(schemaOID),
		"table_name": tbl.name,
	}); err != nil {
		return fmt.Errorf("inserting table: %w", err)
	}

	// TODO: insert indexes.

	return nil
}

// TODO: Return value with status, inserted primary key, etc.
func (tbl *table) Insert(row map[string]interface{}) error {
	remainingColumns := len(row)

	inTuple := NewTupleWithSchema(tbl.schema)
	for idx, elem := range tbl.schema.Elements {
		col := elem.Name
		val, ok := row[col]
		if ok {
			remainingColumns--
		} else {
			// TODO: Enforce null constraints.
			// TODO: Allow a schema element to have a default value.
			val = elem.Type.DefaultValue() // XXX: This is a hack to get some value serialized.
		}
		if err := inTuple.SetDynamic(idx, val); err != nil {
			return err
		}
	}

	if remainingColumns > 0 {
		return errors.New("field found in input that are not present in schema")
	}

	// Key is the primary key prefixed by the table's oid.
	// XXX: This assumes that the 0th element of the tuple is the primary key. This should be
	// configured in schema.
	key := NewTuple(int32(tbl.oid)).Concat(inTuple.Slice(0, 1)).Serialize()
	val, err := tbl.serde.Serialize(inTuple)
	if err != nil {
		return fmt.Errorf("serializing value: %w", err)
	}

	return tbl.db.store.Set(key, val)
}

type scanFunc = func(t *Tuple) bool

func (tbl *table) Scan(fn scanFunc) error {
	tx := tbl.db.store.ReadTransaction()
	defer tx.End()

	prefix := NewTuple(int32(tbl.oid)).Serialize()
	it := tbl.db.store.Scan(tx, ScanArgs{
		Prefix: prefix,
	})
	defer it.Close()

	for it.Next() {
		if entry := it.Item(); entry != nil {
			rowTup, err := tbl.serde.Deserialize(entry.Value)
			if err != nil {
				return fmt.Errorf("decoding row: %w", err)
			}
			if !fn(rowTup) {
				return nil
			}
		}
	}

	return it.Err()
}

// Getting a record:
// 1. Go to the "table" table
// 2. ScanOne for `name = ?`
// 3. If not found, error
// 4. Get oid of table
// 5. Use <oid>.Serialize for get prefix
// (should cache for future reference)
