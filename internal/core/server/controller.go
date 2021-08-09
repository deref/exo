package server

type Controller interface {
	InitResource(spec, state string) error
	MarshalState() (state string, err error)
}
