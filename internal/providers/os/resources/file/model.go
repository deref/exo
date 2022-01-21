package file

import (
	"bufio"
	"context"
	"fmt"
	"io/ioutil"
	"net/textproto"
	"strings"
)

type Model struct {
	HostID   string
	Path     string
	Modified string
	Contents string
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
	contents, _ := ioutil.ReadAll(r)
	m.Contents = string(contents)
	return nil
}

func (m *Model) MarshalModel(ctx context.Context) (string, error) {
	var sb strings.Builder
	// TODO: Validate no linebreaks in header values.
	fmt.Fprintf(&sb, "Host-ID: %s\n", m.HostID)
	fmt.Fprintf(&sb, "Path: %s\n", m.Path)
	fmt.Fprintln(&sb)
	sb.WriteString(m.Contents)
	return sb.String(), nil
}
