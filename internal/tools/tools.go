// NOTE [GO_TOOLS]: In order to ensure we get the correct version of any
// Go-based tools that we depend on, they must be listed in our go.mod and
// go.sum files. However, linters will complain that the tools are unused
// dependencies. This file is ignored during builds (unless you specify the
// "tools" flag), but it is enough to disable the linter warnings.
//
// To install a new tool:
// 1) Add it here.
// 2) Run `go mod tidy`.
// 3) Add it to the `install-tools.sh` script.

//go:build tools

package tools

import (
	_ "github.com/deref/extractgqlts"
)
