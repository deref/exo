package client

import (
	"net/http"
	"net/url"

	josh "github.com/deref/exo/josh/client"
)

type Root struct {
	HTTP *http.Client
	URL  string
}

func (root *Root) Kernel() *Kernel {
	return GetKernel(&josh.Client{
		HTTP: http.DefaultClient,
		URL:  root.URL + "kernel",
	})
}

// TODO: Delete this. Have the kernel return urls as capabilities.
func (root *Root) GetWorkspace(id string) *Workspace {
	return GetWorkspace(&josh.Client{
		HTTP: root.HTTP,
		URL:  root.URL + "workspace?id=" + url.QueryEscape(id),
	})
}
