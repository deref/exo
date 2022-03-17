package resolvers

import (
	"net/http"

	"github.com/deref/exo/internal/util/errutil"
	"github.com/mattn/go-sqlite3"
)

type ErrorResolver struct {
	Message string
}

func isSqlConflict(err error) bool {
	if err, ok := err.(sqlite3.Error); ok {
		return err.ExtendedCode == sqlite3.ErrConstraintUnique
	}
	return false
}

func conflictErrorf(format string, v ...any) error {
	return errutil.HTTPErrorf(http.StatusConflict, format, v...)
}
