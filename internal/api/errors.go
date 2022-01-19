package api

import (
	"fmt"

	"github.com/graph-gophers/graphql-go/errors"
)

type QueryErrorSet []*errors.QueryError

func (errs QueryErrorSet) Error() string {
	switch len(errs) {
	case 0:
		return "<empty QueryErrorSet>"
	case 1:
		return errs[0].Error()
	default:
		return fmt.Sprintf("1st of %d: %v", len(errs), errs[0])
	}
}

func (errs QueryErrorSet) Unwrap() error {
	if len(errs) != 1 {
		return nil
	}
	return errs[0]
}
