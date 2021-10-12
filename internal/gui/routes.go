package gui

import (
	"fmt"
	"net/url"
)

type Routes struct {
	URL string
}

func (r *Routes) NewWorkspaceURL(root string) string {
	return fmt.Sprintf("%s/#/new-project?root=%s", r.URL, url.QueryEscape(root))
}

func (r *Routes) WorkspaceURL(id string) string {
	return fmt.Sprintf("%s/#/workspaces/%s", r.URL, url.PathEscape(id))
}

func (r *Routes) JobURL(id string) string {
	return fmt.Sprintf("%s/#/jobs/%s", r.URL, url.PathEscape(id))
}
