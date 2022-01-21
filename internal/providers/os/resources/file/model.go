package file

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/textproto"
	"strconv"
	"strings"
)

type Model struct {
	HostID  string
	Path    string
	Content string
}

func (m *Model) UnmarshalModel(ctx context.Context, s string) error {
	r := bufio.NewReader(strings.NewReader(s))
	tpr := textproto.NewReader(r)
	hdr, err := tpr.ReadMIMEHeader()
	if err != nil {
		return fmt.Errorf("reading header: %w", err)
	}
	m.HostID = hdr.Get("Host-ID")
	m.Path = hdr.Get("Path")
	content, _ := ioutil.ReadAll(r)

	lengthStr := hdr.Get("Content-Length")
	if lengthStr != "" {
		length, err := strconv.Atoi(lengthStr)
		trimmed := trimEOL(content)
		switch {
		case err != nil:
			return fmt.Errorf("parsing Content-Length: %w", err)
		case len(content) < length:
			return errors.New("content truncated")
		case length == len(trimmed):
			content = trimmed
		case length != len(content):
			return errors.New("content length mismatch")
		}
	}
	m.Content = string(content)
	return nil
}

func (m *Model) MarshalModel(ctx context.Context) (string, error) {
	var sb strings.Builder
	fmt.Fprintf(&sb, "Host-ID: %s\n", m.HostID)
	fmt.Fprintf(&sb, "Path: %s\n", m.Path) // TODO: Validate no line breaks.
	fmt.Fprintf(&sb, "Content-Length: %d\n", len(m.Content))
	fmt.Fprintln(&sb)
	sb.WriteString(m.Content)
	return sb.String(), nil
}

func trimEOL(bs []byte) []byte {
	return bytes.TrimSuffix(bytes.TrimSuffix(bs, []byte("\n")), []byte("\r"))
}
