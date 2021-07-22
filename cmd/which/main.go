package main

import (
	"fmt"
	"os"

	"github.com/deref/exo/util/cmdutil"
	"github.com/deref/exo/util/which"
)

func main() {
	args := os.Args[1:]
	wd, err := os.Getwd()
	if err != nil {
		cmdutil.Fatalf("getting working directory: %w", err)
	}
	pathVar, _ := os.LookupEnv("PATH")
	for _, arg := range args {
		found, err := which.Query{
			WorkingDirectory: wd,
			PathVariable:     pathVar,
			Program:          arg,
		}.Run()
		if err != nil {
			cmdutil.Fatal(err)
		}
		fmt.Println(found)
	}
}
