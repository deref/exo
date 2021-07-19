package server

import "net/http"

type IntrospectionHandler struct {
	MethodNames []string
}

func (h *IntrospectionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
	writeJSON(w, output)
}
