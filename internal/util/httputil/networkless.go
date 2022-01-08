package httputil

import (
	"net/http"
	"net/http/httptest"
)

type NetworklessTransport struct {
	Handler http.Handler
}

func (t *NetworklessTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	recorder := httptest.NewRecorder()
	t.Handler.ServeHTTP(recorder, req)
	res := recorder.Result()
	return res, nil
}
