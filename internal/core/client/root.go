package client

import (
	"net/http"
	"net/url"
	"regexp"

	josh "github.com/deref/exo/internal/josh/client"
)

type Root struct {
	HTTP  *http.Client
	URL   string
	Token string
}

func (root *Root) Kernel() *Kernel {
	return GetKernel(&josh.Client{
		HTTP:  http.DefaultClient,
		URL:   root.URL + "kernel",
		Token: root.Token,
	})
}

// TODO: Delete this. Have the kernel return urls as capabilities.
func (root *Root) GetWorkspace(id string) *Workspace {
	return GetWorkspace(&josh.Client{
		HTTP:  root.HTTP,
		URL:   root.URL + "workspace?id=" + url.QueryEscape(id),
		Token: root.Token,
	})
}

// TODO: Avoid this regex game. Can a synchronous, error-free ID() method be
// baked in to the Workspace API somehow?
func (ws *Workspace) ID() string {
	return regexp.MustCompile("id=(.*)").FindStringSubmatch(ws.client.URL)[1]
}
