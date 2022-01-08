package resolvers

import (
	"context"
	"database/sql"
	"errors"
)

func (r *RootResolver) getRowByID(ctx context.Context, dest interface{}, q string, id *string) error {
	if id == nil {
		return nil
	}
	err := r.DB.GetContext(ctx, dest, q, id)
	if errors.Is(err, sql.ErrNoRows) {
		err = nil
	}
	return err
}
