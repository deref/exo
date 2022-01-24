package resolvers

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/deref/exo/internal/util/errutil"
	"github.com/jmoiron/sqlx"
)

func (r *RootResolver) getRowByKey(ctx context.Context, dest interface{}, q string, id *string) error {
	if id == nil {
		return nil
	}
	err := r.DB.GetContext(ctx, dest, q, id)
	if errors.Is(err, sql.ErrNoRows) {
		err = nil
	}
	return err
}

func trimmed(s string, fallback string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return fallback
	}
	return s
}

func trimmedPtr(p *string, fallback string) *string {
	if p == nil {
		return &fallback
	}
	s := trimmed(*p, fallback)
	return &s
}

func stringPtr(s string) *string {
	return &s
}

func (r *RootResolver) insertRow(ctx context.Context, table string, row interface{}) error {
	v := reflect.ValueOf(row)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	typ := v.Type()
	n := typ.NumField()
	columns := make([]string, 0, n)
	values := make([]interface{}, 0, n)
	for i := 0; i < n; i++ {
		tag := typ.Field(i).Tag.Get("db")
		columns = append(columns, tag)
		values = append(values, v.Field(i).Interface())
	}

	var q strings.Builder
	q.WriteString("INSERT INTO ")
	q.WriteString(table)
	q.WriteString(" (")
	prefix := " "
	for _, column := range columns {
		q.WriteString(prefix)
		q.WriteString(column)
		prefix = ", "
	}
	q.WriteString(" ) VALUES (")
	placeholder := " ?"
	for range columns {
		q.WriteString(placeholder)
		placeholder = ", ?"
	}
	q.WriteString(" )")

	_, err := r.DB.ExecContext(ctx, q.String(), values...)
	return err
}

func transact(ctx context.Context, db *sqlx.DB, f func(tx *sqlx.Tx) error) error {
	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("beginning transaction: %w", err)
	}
	if err := errutil.Recovering(func() error {
		return f(tx)
	}); err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("committing: %q", err)
	}
	return nil
}
