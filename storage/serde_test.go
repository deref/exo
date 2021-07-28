package storage_test

import (
	"testing"

	"github.com/deref/exo/storage"
	"github.com/stretchr/testify/assert"
)

func TestSchematizedRowRoundtrip(t *testing.T) {
	tests := []struct {
		name string
		tup  *storage.Tuple
	}{
		{
			name: "single fixed-length field",
			tup:  storage.NewTuple(int64(1234)),
		},
		{
			name: "multiple fixed-length fields",
			tup: storage.NewTuple(
				int64(1234),
				uint64(3848),
			),
		},
		{
			name: "single variable-length field",
			tup:  storage.NewTuple("Hello World"),
		},
		{
			name: "multiple variable-length fields",
			tup:  storage.NewTuple("Hello World", "Здравей свят"),
		},
		{
			name: "mixed fixed and variable-length fields",
			tup: storage.NewTuple(
				int64(1234),
				"Andrew",
				"Meredith",
				uint64(3848),
				"Author of this test",
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			serde := storage.NewSchematizedRowSerde(tt.tup.Schema())
			data, err := serde.Serialize(tt.tup)
			if !assert.NoError(t, err) {
				return
			}
			deserialized, err := serde.Deserialize(data)
			if !assert.NoError(t, err) {
				return
			}

			assert.Equal(t, tt.tup, deserialized)
		})
	}
}

func BenchmarkSchematizedRowRoundtrip(b *testing.B) {
	tup := storage.NewTuple(
		int64(1234),
		"Andrew",
		"Meredith",
		uint64(3848),
		"Author of this benchmark",
	)
	s := storage.NewSchematizedRowSerde(tup.Schema())
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		data, _ := s.Serialize(tup)
		s.Deserialize(data)
	}
}
