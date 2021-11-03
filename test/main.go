package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"sync"
	"time"

	"github.com/deref/exo/test/tester"
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

func doTest(ctx context.Context, test tester.ExoTest, testName, exoBinPath, fixtureBasePath string, outputMutex *sync.Mutex) error {
	tester := tester.MakeExoTester(exoBinPath, fixtureBasePath, test)

	fmt.Println("running", testName)
	if outputReader, err := tester.RunTest(ctx, test); err != nil {
		outputMutex.Lock()
		defer outputMutex.Unlock()

		output, readErr := ioutil.ReadAll(outputReader)
		if readErr != nil {
			return fmt.Errorf("reading output of failed test: %w", readErr)
		}

		fmt.Printf("test output for %q:\n%s\n", testName, string(output))
		return fmt.Errorf("test failed: %w", err)
	}
	return nil

}

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: e2etest exoBinPath fixtureBasePath")
		os.Exit(1)
	}
	exoBinPath := os.Args[1]
	fixtureBasePath := os.Args[2]

	var outputMutex sync.Mutex

	ctx := context.Background()

	// FIXME: these tests should be parallelisable but there is evidently a race
	// that causes issues.
	for testName := range tests {
		testName := testName
		test := tests[testName]
		exoBinPath := exoBinPath
		fixtureBasePath := fixtureBasePath
		if err := doTest(ctx, test, testName, exoBinPath, fixtureBasePath, &outputMutex); err != nil {
			fmt.Println("Tests failed: ", err)
			os.Exit(1)

		}
	}

}
