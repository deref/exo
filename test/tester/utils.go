package tester

import (
	"context"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"time"
)

// PortIsBound checks if a given TCP port is bound on the local network interface.
func PortIsBound(port int, timeout time.Duration) bool {
	conn, err := net.DialTimeout("tcp", net.JoinHostPort("localhost", strconv.Itoa(port)), timeout)
	if conn != nil {
		defer conn.Close()
	}
	return err == nil
}

func ExpectResponse(ctx context.Context, endpoint string, expectedResponse string) error {
	req, err := http.NewRequestWithContext(ctx, "GET", endpoint, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return fmt.Errorf("failed to get response from %q: %w", endpoint, err)
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("expected status code 200, got %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if string(body) != expectedResponse {
		return fmt.Errorf("expected response body to be %q, got %q", expectedResponse, string(body))
	}
	return nil
}
