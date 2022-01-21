package sdk

import (
	"net/http"

	"github.com/deref/exo/internal/util/errutil"
)

var ErrResourceGone = errutil.NewHTTPError(http.StatusGone, "resource gone")
