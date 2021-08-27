package deps

type StringNode string

func (s StringNode) ID() string {
	return string(s)
}
