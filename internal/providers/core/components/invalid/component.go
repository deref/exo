package invalid

type Invalid struct {
	Err error
}

func (invalid *Invalid) InitResource() error {
	return nil
}

func (invalid *Invalid) MarshalState() (state string, err error) {
	return "", invalid.Err
}

func (invalid *Invalid) IsDeleted() bool {
	return true
}

func (invalid *Invalid) MarkDeleted() {}
