package sdk

import (
	"net/http"

	"github.com/deref/exo/internal/util/errutil"
)

// TODO: Do not abuse HTTP errors.
var ErrResourceGone = errutil.NewHTTPError(http.StatusGone, "resource gone")
var ErrMethodNotAllowed = errutil.NewHTTPError(http.StatusMethodNotAllowed, "method not allowed")
var ErrUnprocessable = errutil.NewHTTPError(http.StatusUnprocessableEntity, "unprocessable entity")
