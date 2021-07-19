package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"

	"github.com/deref/exo/util/jsonutil"
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
	url := c.URL + method
	contentType := "application/json"
	resp, err := c.HTTP.Post(url, contentType, bytes.NewBuffer(inputB))
	if err != nil {
		return fmt.Errorf("posting: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		bs, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("reading response body: %w", err)
		}
		bs = bytes.TrimSpace(bs)
		if len(bs) == 0 {
			return errors.New("empty response")
		}
		return fmt.Errorf("%s", bs)
	}
	outputValue := reflect.ValueOf(output).Elem()
	outputValue.Set(reflect.Zero(outputValue.Type()))
	if err := jsonutil.UnmarshalReader(resp.Body, output); err != nil {
		return fmt.Errorf("unmarshalling output: %w", err)
	}
	return nil
}
