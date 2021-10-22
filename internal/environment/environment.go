package environment

type Source interface {
	EnvironmentSource() string
	ExtendEnvironment(Builder) error
}

type Builder interface {
	AppendVariable(src Source, name string, value string)
}
