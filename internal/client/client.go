package client

import (
	"context"
	"net/http"

	machinebox "github.com/machinebox/graphql"
)

type Client struct {
	gql *machinebox.Client
}

func NewClient(url string, httpClient *http.Client) *Client {
	return &Client{
		gql: machinebox.NewClient(url, machinebox.WithHTTPClient(httpClient)),
	}
}

func (cl *Client) Do(ctx context.Context, doc string, vars map[string]interface{}, res interface{}) error {
	req := machinebox.NewRequest(doc)
	for k, v := range vars {
		req.Var(k, v)
	}
	return cl.gql.Run(ctx, req, res)
}

func (cl *Client) Enqueue(ctx context.Context, mutation string, vars map[string]interface{}) (jobID string, err error) {
	panic("NOT YET IMPLEMENTED")
}

func (cl *Client) Await(ctx context.Context, jobID string) error {
	panic("NOT YET IMPLEMENTED")
}
