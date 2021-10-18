package environment

import (
	"os"
	"strings"
)

type OS struct{}

func (src *OS) EnvironmentSource() string {
	return "server environment"
}

func (src *OS) ExtendEnvironment(b Builder) error {
	// TODO: This should probably somehow shell-out to get the user's current
	// environment, otherwise changes to shell profiles won't take effect until
	// the exo daemon is exited and restarted.
	for _, assign := range os.Environ() {
		parts := strings.SplitN(assign, "=", 2)
		key := parts[0]
		val := parts[1]
		b.AppendVariable(src, key, val)
	}
	return nil
}
