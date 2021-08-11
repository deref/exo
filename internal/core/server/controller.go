package server

type Controller interface {
	InitResource(componentID, spec, state string) error
	MarshalState() (state string, err error)
	IsDeleted() bool
	MarkDeleted()
}
