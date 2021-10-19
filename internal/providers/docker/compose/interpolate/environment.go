package interpolate

type Environment interface {
	Lookup(key string) (value string, found bool)
}

type MapEnvironment map[string]string

func (env MapEnvironment) Lookup(key string) (value string, found bool) {
	value, found = map[string]string(env)[key]
	return
}
