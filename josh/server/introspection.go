package server

import (
	"net/http"

	"github.com/deref/exo/util/httputil"
)

type IntrospectionHandler struct {
	MethodNames []string
}

func (h *IntrospectionHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if handleOptions(w, req) {
		return
	}
	if req.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	type method struct {
		Name string `json:"name"`
	}
	var output struct {
		Methods []method `json:"methods"`
	}
	for _, methodName := range h.MethodNames {
		output.Methods = append(output.Methods, method{
			Name: methodName,
		})
	}
	httputil.WriteJSON(w, http.StatusOK, output)
}
