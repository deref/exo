package client

import (
	"context"
	"net/http"

	machinebox "github.com/machinebox/graphql"
	shurcool "github.com/shurcooL/graphql"
)

type Client struct {
	url        string
	shurcool   *shurcool.Client
	machinebox *machinebox.Client
}

func NewClient(url string, httpClient *http.Client) *Client {
	return &Client{
		url:        url,
		shurcool:   shurcool.NewClient(url, httpClient),
		machinebox: machinebox.NewClient(url, machinebox.WithHTTPClient(httpClient)),
	}
}

// Construct and run a query reflectively from q.
func (cl *Client) Query(ctx context.Context, q interface{}, variables map[string]interface{}) error {
	return cl.shurcool.Query(ctx, q, variables)
}

// Construct and run a mutation reflectively from m.
func (cl *Client) Mutate(ctx context.Context, m interface{}, variables map[string]interface{}) error {
	return cl.shurcool.Mutate(ctx, m, variables)
}

// Run the given query string and decode the response in to resp.
func (cl *Client) Run(ctx context.Context, q string, resp interface{}, vars map[string]interface{}) error {
	req := machinebox.NewRequest(q)
	for k, v := range vars {
		req.Var(k, v)
	}
	return cl.machinebox.Run(ctx, req, resp)
}
