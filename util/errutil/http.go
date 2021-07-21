package errutil

import (
	"fmt"
	"net/http"
)

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
	message string
}

func (err *httpError) Error() string {
	return err.message
}

func (err *httpError) HTTPStatus() int {
	return err.status
}

func NewHTTPError(status int, message string) HTTPError {
	return &httpError{
		status:  status,
		message: message,
	}
}

func HTTPErrorf(status int, format string, v ...interface{}) HTTPError {
	return NewHTTPError(status, fmt.Sprintf(format, v...))
}
