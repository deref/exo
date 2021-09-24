package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"

	"github.com/deref/exo/internal/util/jsonutil"
)

type Client struct {
	HTTP  *http.Client
	URL   string
	Token string
}

func (c *Client) Invoke(ctx context.Context, method string, input interface{}, output interface{}) error {
	inputB, err := json.Marshal(input)
	if err != nil {
		return fmt.Errorf("marshalling input: %w", err)
	}

	endpoint, err := url.Parse(c.URL)
	if err != nil {
		return fmt.Errorf("invalid endpoint: %w", err)
	}
	endpoint.Path += "/" + method

	req, err := http.NewRequest("POST", endpoint.String(), bytes.NewBuffer(inputB))
	if err != nil {
		return fmt.Errorf("forming request: %w", err)
	}
	req.Header.Add("Content-Type", "application/json")
	if c.Token != "" {
		req.Header.Add("Authorization", "Bearer "+c.Token)
	}

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return fmt.Errorf("posting: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		bs, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("reading response body: %w", err)
		}
		var obj struct {
			Message string `json:"message"`
		}
		if resp.Header.Get("content-type") == "application/json" {
			_ = json.Unmarshal(bs, &obj)
		}
		if obj.Message == "" {
			bs = bytes.TrimSpace(bs)
			if len(bs) == 0 {
				return errors.New("empty response")
			}
			return fmt.Errorf("%s", bs)
		}
		return errors.New(obj.Message)
	}
	outputValue := reflect.ValueOf(output).Elem()
	outputValue.Set(reflect.Zero(outputValue.Type()))
	if err := jsonutil.UnmarshalReader(resp.Body, output); err != nil {
		return fmt.Errorf("unmarshalling output: %w", err)
	}
	return nil
}
