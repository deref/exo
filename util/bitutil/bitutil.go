package bitutil

// FlagSetInByte returns true iff the `pos` bit is set to true in `flags`.
func FlagSetInByte(flags byte, flag byte) bool {
	return (flags>>flag)&1 == 1
}
