package httputil

import (
	_ "embed"
	"errors"
	"fmt"
	"log"
	"net/http"
	"text/template"

	"github.com/deref/exo/util/errutil"
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
	if status == http.StatusInternalServerError {
		log.Printf("error processing request: %v", err)
	}

	contentType, _ := accept.Negotiate(req.Header.Get("accept"), "application/json", "text/html", "text/plain")
	switch contentType {
	case "application/json":
		WriteJSON(w, status, map[string]interface{}{
			"status":  status,
			"message": message,
		})
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
