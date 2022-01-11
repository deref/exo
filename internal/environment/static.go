package environment

type Static struct {
	Name      string
	Variables map[string]string
}

func (src *Static) EnvironmentSource() string {
	return src.Name
}

func (src *Static) ExtendEnvironment(b Builder) error {
	for k, v := range src.Variables {
		b.AppendVariable(src, k, v)
	}
	return nil
}

var Default = &Static{
	Name: "exo",
	Variables: map[string]string{
		// Encourage programs to log with colors enabled.  The closest thing to a
		// standard for this is <https://bixense.com/clicolors/>, but
		// support is spotty. This may grow if there are other popular enviornment
		// variables to include. If we grow PTY support, this may become unnecessary.
		"CLICOLOR":       "1",
		"CLICOLOR_FORCE": "1",
		"FORCE_COLOR":    "3", // https://github.com/chalk/chalk/tree/9d5b9a133c3f8aa9f24de283660de3f732964aaa#supportscolor
	},
}
