package environment

import (
	"github.com/deref/exo/internal/esv"
)

type ESV struct {
	Client esv.EsvClient
	Name   string
	URL    string
}

func (src *ESV) EnvironmentSource() string {
	return src.Name
}

func (src *ESV) ExtendEnvironment(b Builder) error {
	secrets, err := src.Client.GetWorkspaceSecrets(src.URL)
	if err != nil {
		return err
	}
	for k, v := range secrets {
		b.AppendVariable(src, k, v)
	}
	return nil
}
