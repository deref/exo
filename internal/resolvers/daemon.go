package resolvers

import (
	"os"
	"time"

	"github.com/deref/exo/internal/util/osutil"
)

func (r *MutationResolver) StopDaemon() *Void {
	// Shutdown asynchronously and make a best effort to return synchronously.
	// TODO: Reliably acknowledge the exit request before exiting.
	go func() {
		ownPid := os.Getpid()
		_ = osutil.TerminateProcessWithTimeout(ownPid, 5*time.Second)
	}()
	return &Void{}
}
