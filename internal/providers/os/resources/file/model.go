package file

import (
	"bufio"
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

	length, err := strconv.Atoi(hdr.Get("Content-Length"))
	switch {
	case err != nil:
		// Ignore invalid length.
	case len(content) < length:
		return errors.New("content truncated")
	case length < len(content):
		// Primarily to chop off a trailing newline.
		content = content[:length]
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
