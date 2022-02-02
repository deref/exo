package main

import (
	"fmt"
	"os"

	"github.com/deref/exo/internal/util/cmdutil"
	"github.com/deref/exo/internal/util/which"
)

func main() {
	args := os.Args[1:]
	for _, arg := range args {
		found, err := which.Which(arg)
		if err != nil {
			cmdutil.Fatal(err)
		}
		fmt.Println(found)
	}
}
