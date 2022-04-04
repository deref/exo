package sdk

import (
	"net/http"

	"github.com/deref/exo/internal/util/errutil"
)

// TODO: Do not abuse HTTP errors.

// Resource controllers use this to signal that an external resource has been deleted.
var ErrResourceGone = errutil.NewHTTPError(http.StatusGone, "resource gone")

// Resource controllers use this to signal that a resource update is not
// possible and instead the resource must be recreated.
func NewNotImplementedErrorf(format string, args ...interface{}) errutil.HTTPError {
	return errutil.HTTPErrorf(http.StatusNotImplemented, "not implemented: "+format, args...)
}
