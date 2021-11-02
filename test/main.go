package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/deref/exo/test/tester"
	"golang.org/x/sync/errgroup"
)

type exoTestParams struct {
	exoBinary string
}

var basicT0Test = func(ctx context.Context, t tester.ExoTester) error {
	if _, _, err := t.RunExo(ctx, "init"); err != nil {
		return err
	}
	if _, _, err := t.RunExo(ctx, "start"); err != nil {
		return err
	}
	return t.WaitTillProcessRunning(ctx, "t0", time.Second*5)
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
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: e2etest exoBinPath fixtureBasePath")
		os.Exit(1)
	}
	exoBinPath := os.Args[1]
	fixtureBasePath := os.Args[2]

	eg, ctx := errgroup.WithContext(context.Background())
	for testName, test := range tests {
		eg.Go(func() error {
			tester := tester.MakeExoTester(exoBinPath, fixtureBasePath, test)

			fmt.Println("running", testName)
			if err := tester.RunTest(ctx, test); err != nil {
				return fmt.Errorf("test failed: %w", err)
			}
			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		fmt.Println("Tests failed: ", err)
		os.Exit(1)
	}
}
