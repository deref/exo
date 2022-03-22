package cmdutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatBytes(t *testing.T) {
	// Bytes are presented exactly.
	assert.Equal(t, "0 B", FormatBytes(0))
	assert.Equal(t, "512 B", FormatBytes(512))
	// All others are presented with decimals.
	assert.Equal(t, "1 KiB", FormatBytes(1024))
	assert.Equal(t, "1.5 KiB", FormatBytes(1536))
	assert.Equal(t, "1.75 KiB", FormatBytes(1792))
	assert.Equal(t, "1.75 KiB", FormatBytes(1792))
	// Max precision of two decimal places.
	assert.Equal(t, "1.07 KiB", FormatBytes(1100)) // Rounding down.
	assert.Equal(t, "1.08 KiB", FormatBytes(1101)) // Rounding up.
	// Up to three digits left of the decimal.
	assert.Equal(t, "123 MiB", FormatBytes(128974848))
	// Largest unit is EiB.
	assert.Equal(t, "1 EiB", FormatBytes(1152921504606846976))
	assert.Equal(t, "16 EiB", FormatBytes(18446744073709551615))
}
