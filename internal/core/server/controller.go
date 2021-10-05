package server

type Controller interface {
	InitResource() error
	MarshalState() (state string, err error)
}
