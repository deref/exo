package errutil

import (
	"errors"
	"fmt"
	"net/http"
)

// TODO: Stop abusing this for internal vs external errors.
type HTTPError interface {
	error
	HTTPStatus() int
}

func HTTPStatus(err error) int {
	if err, ok := err.(HTTPError); ok {
		return err.HTTPStatus()
	}
	return http.StatusInternalServerError
}

type httpError struct {
	status  int
	wrapped error
}

func (err *httpError) HTTPStatus() int {
	return err.status
}

func (err *httpError) Error() string {
	return err.wrapped.Error()
}

func (err *httpError) Unwrap() error {
	return err.wrapped
}

func WithHTTPStatus(status int, err error) HTTPError {
	return &httpError{
		status:  status,
		wrapped: err,
	}
}

func NewHTTPError(status int, message string) HTTPError {
	return WithHTTPStatus(status, errors.New(message))
}

func HTTPErrorf(status int, format string, v ...any) HTTPError {
	return WithHTTPStatus(status, fmt.Errorf(format, v...))
}

var InternalServerError = NewHTTPError(http.StatusInternalServerError, "internal server error")
