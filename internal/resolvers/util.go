package resolvers

import (
	"context"
	"database/sql"
	"errors"
	"strings"
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
