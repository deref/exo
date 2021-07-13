package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"path"

	"github.com/deref/exo/jsonutil"
)

type Client struct {
	HTTP *http.Client
	URL  string
}

func (c *Client) Invoke(ctx context.Context, method string, input interface{}, output interface{}) error {
	inputB, err := json.Marshal(input)
	if err != nil {
		return fmt.Errorf("marshalling input: %w", err)
	}
	url := path.Join(c.URL, method)
	contentType := "application/json"
	resp, err := c.HTTP.Post(url, contentType, bytes.NewBuffer(inputB))
	if err != nil {
		return fmt.Errorf("posting: %w", err)
	}
	if err := jsonutil.UnmarshalReader(resp.Body, output); err != nil {
		return fmt.Errorf("unmarshalling output: %w", err)
	}
	return nil
}
