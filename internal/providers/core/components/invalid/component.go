package invalid

type Invalid struct {
	Err error
}

func (invalid *Invalid) InitResource() error {
	return invalid.Err
}

func (invalid *Invalid) MarshalState() (state string, err error) {
	return "", invalid.Err
}
