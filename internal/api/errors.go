package api

import (
	"fmt"
	"strings"

	"github.com/graph-gophers/graphql-go/errors"
)

type QueryErrorSet []*errors.QueryError

func (errs QueryErrorSet) Error() string {
	switch len(errs) {
	case 0:
		return "<empty QueryErrorSet>"
	case 1:
		return queryErrorString(errs[0])
	default:
		return fmt.Sprintf("1st of %d: %s", len(errs), queryErrorString(errs[0]))
	}
}

// Mimics errors.QueryError.Error(), but without the "graphql: " prefix.
func queryErrorString(err *errors.QueryError) string {
	if err == nil {
		return "<nil>"
	}
	var b strings.Builder
	b.WriteString(err.Message)
	for _, loc := range err.Locations {
		fmt.Fprintf(&b, " (line %d, column %d)", loc.Line, loc.Column)
	}
	return b.String()
}

func (errs QueryErrorSet) Unwrap() error {
	if len(errs) != 1 {
		return nil
	}
	return errs[0]
}
