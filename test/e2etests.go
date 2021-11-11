package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/deref/exo/test/tester"
)

var basicT0Test = func(ctx context.Context, t tester.ExoTester) error {
	if _, _, err := t.RunExo(ctx, "init"); err != nil {
		return err
	}
	if _, _, err := t.RunExo(ctx, "start"); err != nil {
		return err
	}
	ctx, cancel := context.WithDeadline(ctx, time.Now().Add(time.Second*5))
	defer cancel()
	return t.WaitTillProcessesRunning(ctx, "t0")
}

var tests = map[string]tester.ExoTest{
	"basic-procfile": {
		FixtureDir: "basic-procfile",
		Test:       basicT0Test,
	},
	"basic-dockerfile": {
		FixtureDir: "basic-dockerfile",
		Test:       basicT0Test,
	},
	"basic-exo-hcl": {
		FixtureDir: "basic-exo-hcl",
		Test:       basicT0Test,
	},
	"simple-example": {
		FixtureDir: "simple-example",
		Test: func(ctx context.Context, t tester.ExoTester) error {
			if _, _, err := t.RunExo(ctx, "init"); err != nil {
				return err
			}
			if _, _, err := t.RunExo(ctx, "start"); err != nil {
				return err
			}
			ctx, cancel := context.WithDeadline(ctx, time.Now().Add(time.Second*10))
			defer cancel()
			if err := t.WaitTillProcessesRunning(ctx, "web", "echo", "echo-short"); err != nil {
				return err
			}

			for port := 44222; port <= 44224; port++ {
				timeout := time.Second
				conn, err := net.DialTimeout("tcp", net.JoinHostPort("localhost", strconv.Itoa(port)), timeout)
				if err != nil {
					return fmt.Errorf("failed to connect to port %d: %w", port, err)
				}
				if conn != nil {
					defer conn.Close()
				}
			}

			resp, err := http.Get("http://localhost:44224")
			if err != nil {
				return fmt.Errorf("failed to get response from port 44224: %w", err)
			}

			if resp.StatusCode != 200 {
				return fmt.Errorf("expected status code 200, got %d", resp.StatusCode)
			}

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return fmt.Errorf("failed to read response body: %w", err)
			}

			expectedResponse := "Hi!"
			if string(body) != expectedResponse {
				return fmt.Errorf("expected response body to be %q, got %q", expectedResponse, string(body))
			}

			return nil
		},
	},
}
