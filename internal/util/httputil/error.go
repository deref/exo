package httputil

import (
	_ "embed"
	"errors"
	"fmt"
	"net/http"
	"text/template"

	"github.com/deref/exo/internal/util/errutil"
	"github.com/deref/exo/internal/util/jsonutil"
	"github.com/deref/exo/internal/util/logging"
	accept "github.com/timewasted/go-accept-headers"
)

func WriteError(w http.ResponseWriter, req *http.Request, err error) {
	status := http.StatusInternalServerError
	message := "internal server error"
	var httpErr errutil.HTTPError
	if errors.As(err, &httpErr) {
		status = httpErr.HTTPStatus()
		message = err.Error()
	}
	logger := logging.CurrentLogger(req.Context())
	if status == http.StatusInternalServerError {
		logger.Infof("error processing request: %v", err)
	}

	contentType, _ := accept.Negotiate(req.Header.Get("accept"), "application/json", "text/event-stream", "text/html", "text/plain")
	switch contentType {
	case "application/json":
		WriteJSON(w, req, status, map[string]any{
			"status":  status,
			"message": message,
		})
	case "text/event-stream":
		sse := StartSSE(w)
		// An event of type "error" isn't part of any relevant text/event-stream
		// sub-protocol, but at least somethig will show in the browser. This was
		// added specifically for GraphQL clients which swallow the HTTP status
		// code when the response content-type does not match.
		sse.SendEvent("error", jsonutil.MustMarshal(map[string]any{
			"status":  status,
			"message": message,
		}))

	case "text/html":
		tmpl, parseErr := template.New("error").Parse(errorTemplate)
		if parseErr != nil {
			panic(parseErr)
		}
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(status)
		_ = tmpl.Execute(w, errorData{
			Status:  status,
			Message: err.Error(),
		})
	case "text/plain", "":
		WriteString(w, status, fmt.Sprintf("error %d: %s", status, message))
	}
}

//go:embed error.html.tmpl
var errorTemplate string

type errorData struct {
	Status  int
	Message string
}
