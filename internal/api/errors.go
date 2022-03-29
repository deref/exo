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

// Same information as errors.QueryError.Error(), but formatted for better
// error message composition.
func queryErrorString(err *errors.QueryError) string {
	if err == nil {
		return "<nil>"
	}

	var b strings.Builder

	for _, loc := range err.Locations {
		fmt.Fprintf(&b, "on line %d, column %d: ", loc.Line, loc.Column)
		break
	}

	if len(err.Path) > 0 {
		b.WriteString("resolving")
		sep := " "
		for _, elem := range err.Path {
			fmt.Fprintf(&b, "%s%v", sep, elem)
			sep = "."
		}
		b.WriteString(": ")
	}

	b.WriteString(err.Message)

	return b.String()
}

func (errs QueryErrorSet) Unwrap() error {
	if len(errs) != 1 {
		return nil
	}
	return errs[0]
}

func ToQueryErrorSet(err error) QueryErrorSet {
	switch err := err.(type) {
	case QueryErrorSet:
		return err
	default:
		return QueryErrorSet{ToQueryError(err)}
	}
}

func ToQueryError(err error) *errors.QueryError {
	switch err := err.(type) {
	case *errors.QueryError:
		return err
	default:
		return &errors.QueryError{
			Err:     err,
			Message: err.Error(),
		}
	}
}
