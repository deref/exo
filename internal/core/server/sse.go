package server

import (
	"bytes"
	"net/http"
)

// https://developer.mozilla.org/en-US/docs/Web/API/Server-sent_events/Using_server-sent_events
type SSEWriter struct {
	response http.ResponseWriter
	flusher  http.Flusher
}

func StartSSE(w http.ResponseWriter) *SSEWriter {
	h := w.Header()
	h.Set("Cache-Control", "no-store")
	h.Set("Connection", "keep-alive")
	h.Set("Content-Type", "text/event-stream")
	return &SSEWriter{
		response: w,
		flusher:  w.(http.Flusher),
	}
}

func (w *SSEWriter) KeepAlive() {
	// This could (should?) have an additional trailing newline, but that causes
	// @graphql-sse/client to choke due to the following bug:
	// See https://github.com/Azure/fetch-event-source/issues/15
	w.sendMessage([]byte(":"))
}

// Type and data may not contain newlines.
func (w *SSEWriter) SendEvent(typ string, data []byte) {
	var buf bytes.Buffer
	if typ != "" {
		buf.WriteString("event: ")
		buf.WriteString(typ)
		buf.WriteString("\n")
	}
	if len(data) > 0 {
		buf.WriteString("data: ")
		buf.Write(data)
		buf.WriteString("\n")
	}
	w.sendMessage(buf.Bytes())
}

func (w *SSEWriter) sendMessage(message []byte) {
	w.response.Write(message)
	w.response.Write([]byte("\n"))
	w.flusher.Flush()
}
