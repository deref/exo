package server

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/deref/exo/internal/api"
	"github.com/deref/exo/internal/util/httputil"
	"github.com/deref/exo/internal/util/jsonutil"
)

// Implements the standard Graphql-over-HTTP spec with support for the
// graphql-sse protocol in "Distinct connections mode".
// https://github.com/graphql/graphql-over-http/blob/main/spec/GraphQLOverHTTP.md
// https://github.com/enisdenjo/graphql-sse/blob/58ac322d8a56e5f4a376bd165ae03c7f76716a03/PROTOCOL.md
type GraphqlHandler struct {
	Service api.Service
}

func (h *GraphqlHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var params struct {
		Query         string         `json:"query"`
		OperationName string         `json:"operationName"`
		Variables     map[string]any `json:"variables"`
	}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var sse *httputil.SSEWriter
	if strings.Contains(r.Header.Get("Accept"), "text/event-stream") {
		sse = httputil.StartSSE(w)
	} else {
		w.Header().Set("Content-Type", "application/json")
	}

	newRes := func() any {
		var res any
		return &res
	}
	// TODO: Make use of params.OperationName.
	sub := h.Service.Subscribe(ctx, newRes, params.Query, params.Variables)
	defer sub.Stop()
	for {
		select {
		case <-time.After(3 * time.Second):
			if sse != nil {
				sse.KeepAlive()
			}
		case event, ok := <-sub.Events():
			var result struct {
				Errors api.QueryErrorSet `json:"errors,omitempty"`
				Data   any               `json:"data,omitempty"`
			}
			if ok {
				result.Data = event
			} else if err := sub.Err(); err != nil {
				// SEE NOTE [GRAPHQL_PARTIAL_FAILURE].
				result.Errors = api.ToQueryErrorSet(err)
			} else {
				if sse != nil {
					sse.SendEvent("complete", nil)
				}
				return
			}
			resultJSON := jsonutil.MustMarshal(result)
			if sse == nil {
				w.Write(resultJSON)
			} else {
				sse.SendEvent("next", resultJSON)
			}
			if !ok {
				return
			}
		}
	}
}
