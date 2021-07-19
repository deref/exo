package gensym

import (
	"bytes"
	"crypto/rand"
	"encoding/base32"
)

var encoding32 = base32.NewEncoding("0123456789abcdefghjkmnpqrstvwxyz").WithPadding(base32.NoPadding)

func RandomBase32() string {
	var bs [16]byte
	n, err := rand.Read(bs[:])
	if n != 16 {
		panic("unexpected end of rand source")
	}
	if err != nil {
		panic(err)
	}
	var buf bytes.Buffer
	enc := base32.NewEncoder(encoding32, &buf)
	_, _ = enc.Write(bs[:])
	_ = enc.Close()
	return buf.String()
}
