package cmdutil

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"text/tabwriter"
)

type TableWriter struct {
	tabWriter   *tabwriter.Writer
	columns     []string
	wroteHeader bool
	buf         bytes.Buffer
}

func NewTableWriter(columns ...string) *TableWriter {
	return &TableWriter{
		tabWriter: tabwriter.NewWriter(os.Stdout, 4, 8, 3, ' ', 0),
		columns:   columns,
	}
}

func (w *TableWriter) WriteRow(values ...string) {
	if !w.wroteHeader {
		w.writeRow("# ", w.columns)
		w.wroteHeader = true
	}
	w.writeRow("", values)
}

func (w *TableWriter) writeRow(prefix string, values []string) {
	n := len(w.columns)
	if len(values) != n {
		panic(fmt.Errorf("expected %d columns, got %d", n, len(values)))
	}
	w.buf.Reset()
	io.WriteString(&w.buf, prefix)
	for i, v := range values {
		io.WriteString(&w.buf, v)
		if i == n-1 {
			io.WriteString(&w.buf, "\n")
		} else {
			io.WriteString(&w.buf, "\t")
		}
	}
	w.tabWriter.Write(w.buf.Bytes())
}

func (w *TableWriter) Flush() {
	w.tabWriter.Flush()
}
