// Log collector in peer mode. Can test stateless behavior normally.
// Can also POST /collect to test the stateful behavior.

package main

import (
	"github.com/deref/exo/logcol"
	"github.com/deref/exo/logcol/api"
	"github.com/deref/pier"
)

func main() {
	lc := logcol.NewLogCollector()
	pier.Main(api.NewLogCollectorMux("/", lc))
}
