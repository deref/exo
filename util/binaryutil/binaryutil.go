package binaryutil

import "errors"

// FlagSetInByte returns true iff the `pos` bit is set to true in `flags`.
func FlagSetInByte(flags byte, pos byte) bool {
	return (flags>>pos)&1 == 1
}

// IncrementBytes returns a byte slice that is incremented by 1 bit.
// If `val` is not already only 255-valued bytes, then it is mutated and returned.
// Otherwise, a new slice is allocated and returned.
func IncrementBytes(val []byte) []byte {
	if val == nil {
		return nil
	}

	for idx := len(val) - 1; idx >= 0; idx-- {
		byt := val[idx]
		if byt == 255 {
			val[idx] = 0
		} else {
			val[idx] = byt + 1
			return val
		}
	}

	// Still carrying from previously most significant byte, so add a new 1-valued byte.
	newVal := make([]byte, len(val)+1)
	newVal[0] = 1
	return newVal
}

// DecrementBytes returns a byte slice that is decremented by 1 bit.
// If `val` is already 0-valued, then an error is returned. Otherwise, `val`
// is mutated.
func DecrementBytes(val []byte) error {
	if val == nil {
		return nil
	}

	for idx := len(val) - 1; idx >= 0; idx-- {
		byt := val[idx]
		if byt == 0 {
			val[idx] = 255
		} else {
			val[idx] = byt - 1
			return nil
		}
	}

	// All bytes in slice were 0-valued.
	return errors.New("already 0")
}
