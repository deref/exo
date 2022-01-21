package errutil

import (
	"bytes"
	"errors"
	"fmt"
	"runtime/debug"
)

// Implementation of the error interface for any panic value.
type Panic struct {
	Value interface{}
}

func (err Panic) Error() string {
	return fmt.Sprintf("panic with non-error value: %v", err.Value)
}

// Coerces non-nil values to errors, wrapping non-errors in a Panic structure.
func ToError(x interface{}) error {
	switch x := x.(type) {
	case nil:
		return nil
	case error:
		return x
	default:
		return Panic{x}
	}
}

// Pair a wrapped error with a stack trace.
type TracedError struct {
	error
	stack string
}

func NewTracedError(err error, stack string) TracedError {
	return TracedError{
		error: err,
		stack: stack,
	}
}

func (err TracedError) Unwrap() error {
	return err.error
}

func (err TracedError) Stack() string {
	return err.stack
}

// Like ToError, but wraps non-nil errors as a TracedError.
func ToTracedError(r interface{}) error {
	switch err := ToError(r).(type) {
	case nil:
		return nil
	case TracedError:
		return err
	default:
		return NewTracedError(err, string(debug.Stack()))
	}
}

// Encapsulates the pattern of writing to a named err result value if recovered.
// Usage: `defer errutil.RecoverTo(&err)`.
func RecoverTo(err *error) {
	if r := ToTracedError(recover()); r != nil {
		*err = r
	}
}

// Invokes f, recovering from any panics.
func Recovering(f func() error) (err error) {
	defer RecoverTo(&err)
	return f()
}

func ErrorWithStack(err error) string {
	var buf bytes.Buffer
	_, _ = buf.WriteString(err.Error())
	var traced TracedError
	if errors.As(err, &traced) {
		_, _ = buf.WriteString("\n")
		_, _ = buf.WriteString(traced.Stack())
	}
	return buf.String()
}

func Coalesce(errors ...error) error {
	for _, err := range errors {
		if err != nil {
			return err
		}
	}
	return nil
}
