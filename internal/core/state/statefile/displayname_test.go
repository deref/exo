package statefile

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDisplayNameBuilder(t *testing.T) {
	filepathSeparator := "/"
	b := newDisplayNameBuilder(filepathSeparator)
	b.AddPath("/personal/unique")
	b.AddPath("/personal/duplicate")
	b.AddPath("/work/duplicate")
	// b.dump()
	assert.Equal(t, "unique", b.GetDisplayName("/personal/unique"))
	assert.Equal(t, "personal/duplicate", b.GetDisplayName("/personal/duplicate"))
	assert.Equal(t, "work/duplicate", b.GetDisplayName("/work/duplicate"))
	assert.Equal(t, "unexpected", b.GetDisplayName("unexpected"))
}
