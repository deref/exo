package gui

import (
	"fmt"
	"net/url"
)

type Routes struct {
	URL string
}

func (r *Routes) NewWorkspaceURL(root string) string {
	return fmt.Sprintf("%s/#/new-workspace?root=%s", r.URL, url.QueryEscape(root))
}

func (r *Routes) WorkspaceURL(id string) string {
	return fmt.Sprintf("%s/#/workspaces/%s", r.URL, url.PathEscape(id))
}
