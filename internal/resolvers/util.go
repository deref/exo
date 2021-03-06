package resolvers

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"cuelang.org/go/cue"
	"github.com/deref/exo/internal/util/cueutil"
	"github.com/deref/exo/internal/util/errutil"
	"github.com/jmoiron/sqlx"
)

func (r *RootResolver) getRowByKey(ctx context.Context, dest any, q string, id *string) error {
	if id == nil {
		return nil
	}
	err := r.db.GetContext(ctx, dest, q, id)
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

// See insertRowEx.
func (r *RootResolver) insertRow(ctx context.Context, table string, row any) error {
	return r.insertRowEx(ctx, table, row, "")
}

// If row is a pointer, it will be updated with the results of `RETURNING *`.
func (r *RootResolver) insertRowEx(ctx context.Context, table string, row any, extra string) error {
	v := reflect.ValueOf(row)
	returning := false
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
		returning = true
	}

	typ := v.Type()
	n := typ.NumField()
	columns := make([]string, 0, n)
	values := make([]any, 0, n)
	for i := 0; i < n; i++ {
		tag := typ.Field(i).Tag.Get("db")
		columns = append(columns, tag)
		values = append(values, v.Field(i).Interface())
	}

	var b strings.Builder
	b.WriteString("INSERT INTO ")
	b.WriteString(table)
	b.WriteString(" (")
	prefix := " "
	for _, column := range columns {
		b.WriteString(prefix)
		b.WriteString(column)
		prefix = ", "
	}
	b.WriteString(" ) VALUES (")
	placeholder := " ?"
	for range columns {
		b.WriteString(placeholder)
		placeholder = ", ?"
	}
	b.WriteString(" )")

	b.WriteString(extra)

	if returning {
		b.WriteString(" RETURNING *")
	}

	q := b.String()

	if returning {
		return r.db.GetContext(ctx, row, q, values...)
	} else {
		_, err := r.db.ExecContext(ctx, q, values...)
		return err
	}
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

func mustSqlIn(query string, args ...any) (string, []any) {
	var err error
	query, args, err = sqlx.In(query, args...)
	if err != nil {
		panic(err)
	}
	return query, args
}

func rowsAffected(res sql.Result) int64 {
	n, err := res.RowsAffected()
	if err != nil {
		// Sqlite supports this, so should not fail.
		panic(err)
	}
	return n
}

func isTrue(b *bool) bool {
	if b == nil {
		return false
	}
	return *b
}

func formatConfiguration(v cue.Value, final bool) (string, error) {
	opts := []cue.Option{
		cue.ResolveReferences(true),
	}
	if final {
		opts = append(opts, cue.Final())
	}
	return cueutil.ValueToString(v, opts...)
}

func isValid(v cue.Value) bool {
	// TODO: Any options to use here?
	return v.Validate() == nil
}

func validateResolve[T any](label string, ref string, resolver *T, err error) error {
	if err != nil {
		return fmt.Errorf("resolving %s: %w", label, err)
	}
	if resolver == nil {
		return errutil.HTTPErrorf(http.StatusNotFound, "no such %s: %q", label, ref)
	}
	return nil
}
