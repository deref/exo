package compose

// Some sections and subsections can be expressed as either a map or sequence.
// Style records the original kind of syntax used in the source YAML.
type Style byte

const (
	UnknownStyle Style = 0
	MapStyle           = 'M'
	SeqStyle           = 'S'
)
