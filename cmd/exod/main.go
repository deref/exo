// Server unbundled from the CLI for development.
// The production entrypoint is the `exo server` command.

package main

import "github.com/deref/exo/exod"

func main() {
	exod.Main()
}
