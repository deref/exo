package invalid

type Invalid struct {
	Err error
}

func (invalid *Invalid) InitResource(componentID, spec, state string) error {
	return nil
}

func (invalid *Invalid) MarshalState() (state string, err error) {
	return "", invalid.Err
}
