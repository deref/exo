package environment

import (
	"github.com/joho/godotenv"
)

type Dotenv struct {
	Path string
}

func (src *Dotenv) EnvironmentSource() string {
	return "env file"
}

func (src *Dotenv) ExtendEnvironment(b Builder) error {
	dotEnvMap, err := godotenv.Read(src.Path)
	if err != nil {
		return err
	}
	for name, value := range dotEnvMap {
		b.AppendVariable(src, name, value)
	}
	return nil
}
