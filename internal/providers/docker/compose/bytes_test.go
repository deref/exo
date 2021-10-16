package compose

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseBytes(t *testing.T) {
	check := func(s string, expected int64) {
		var actual Bytes
		err := actual.Parse(s)
		if assert.NoError(t, err) {
			assert.Equal(t, expected, actual.Int64())
		}
	}
	check("5", 5)
	check("10b", 10)
	check("2k", 2048)
	check("2kb", 2048)
	check("1m", 1024*1024)
	check("1mb", 1024*1024)
	check("1g", 1024*1024*1024)
	check("1gb", 1024*1024*1024)
}

func TestBytesYAML(t *testing.T) {
	testYAML(t, "int", `1234`, Bytes{
		String:   MakeInt(1234).String,
		Quantity: 1234,
	})
	testYAML(t, "string", `5k`, Bytes{
		String:   MakeString("5k"),
		Quantity: 5,
		Unit: ByteUnit{
			Suffix: "k",
			Scalar: 1024,
		},
	})
}
