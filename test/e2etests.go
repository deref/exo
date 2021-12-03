package main

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
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
	return t.WaitTillProcessesReachState(ctx, "running", []string{"t0"})
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

			if err := t.WaitTillProcessesReachState(ctx, "running", []string{"web", "echo", "echo-short"}); err != nil {
				return err
			}

			timeout := time.Second
			for port := 44222; port <= 44224; port++ {
				if !tester.PortIsBound(port, timeout) {
					return fmt.Errorf("port %d is not bound", port)
				}
			}

			if err := tester.ExpectResponse(ctx, "http://localhost:44224", "Hi!"); err != nil {
				return err
			}

			// Check that we can stop the workspace.
			if _, _, err := t.RunExo(ctx, "stop"); err != nil {
				return err
			}
			if err := t.WaitTillProcessesReachState(ctx, "stopped", []string{"web", "echo", "echo-short"}); err != nil {
				return err
			}
			for port := 44222; port <= 44224; port++ {
				if tester.PortIsBound(port, timeout) {
					return fmt.Errorf("port %d is still bound", port)
				}
			}

			// Check that we can start just one process.
			if _, _, err := t.RunExo(ctx, "start", "echo"); err != nil {
				return err
			}
			if err := t.WaitTillProcessesReachState(ctx, "running", []string{"echo"}); err != nil {
				return err
			}
			if err := t.WaitTillProcessesReachState(ctx, "stopped", []string{"web", "echo-short"}); err != nil {
				return err
			}

			return nil
		},
	},
	"api-gateway": {
		FixtureDir: "api-gateway",
		Test: func(ctx context.Context, t tester.ExoTester) error {
			if _, _, err := t.RunExo(ctx, "init"); err != nil {
				return err
			}
			if _, _, err := t.RunExo(ctx, "start"); err != nil {
				return err
			}

			ctx, cancel := context.WithDeadline(ctx, time.Now().Add(time.Second*10))
			defer cancel()

			if err := t.WaitTillProcessesReachState(ctx, "running", []string{"t0", "t1"}); err != nil {
				return err
			}

			validateResponse := func(url string, expectedBody string) error {
				for {
					resp, err := http.Get(url)
					if err != nil {
						return err
					}
					if resp.StatusCode == 200 {
						body, err := ioutil.ReadAll(resp.Body)
						if err != nil {
							return err
						}
						if string(bytes.TrimSpace(body)) != expectedBody {
							return fmt.Errorf("unexpected response %s", body)
						}
						return nil
					}
					select {
					case <-ctx.Done():
						return ctx.Err()
					default:
					}
				}
			}

			if err := validateResponse("http://localhost:44544", "a"); err != nil {
				return err
			}
			if err := validateResponse("http://t1.exo.localhost:44544", "b"); err != nil {
				return err
			}

			return nil
		},
	},
}
