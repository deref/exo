package environment

import (
	"os"

	"github.com/deref/exo/internal/util/osutil"
)

type OS struct{}

func (src *OS) EnvironmentSource() string {
	return "server environment"
}

func (src *OS) ExtendEnvironment(b Builder) error {
	// TODO: This should probably somehow shell-out to get the user's current
	// environment, otherwise changes to shell profiles won't take effect until
	// the exo daemon is exited and restarted.
	for _, entry := range os.Environ() {
		name, value := osutil.ParseEnvvEntry(entry)
		b.AppendVariable(src, name, value)
	}
	return nil
}
