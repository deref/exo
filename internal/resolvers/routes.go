package resolvers

import (
	"fmt"
	"net/url"
	"strings"
)

type RoutesResolver struct {
	Root string
}

func (r *QueryResolver) Routes() *RoutesResolver {
	return &RoutesResolver{
		Root: r.GUIEndpoint,
	}
}

func (r *RoutesResolver) NewProjectURL(args struct {
	Workspace *string
}) string {
	var b strings.Builder
	b.WriteString(r.Root)
	b.WriteString("/#/new-project")
	if args.Workspace != nil {
		b.WriteString("?workspace=")
		b.WriteString(url.QueryEscape(*args.Workspace))
	}
	return b.String()
}

func (r *RoutesResolver) WorkspaceURL(args struct {
	ID string
}) string {
	return r.workspaceURL(args.ID)
}

func (r *RoutesResolver) workspaceURL(id string) string {
	return fmt.Sprintf("%s/#/workspaces/%s", r.Root, url.PathEscape(id))
}

func (r *RoutesResolver) JobURL(args struct {
	ID string
}) string {
	return r.jobURL(args.ID)
}

func (r *RoutesResolver) jobURL(id string) string {
	return fmt.Sprintf("%s/#/jobs/%s", r.Root, url.PathEscape(id))
}
