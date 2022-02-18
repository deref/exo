package hashutil

import (
	"crypto/sha256"
	"fmt"
)

func Sha256Hex(data []byte) string {
	hash := sha256.Sum256(data)
	return fmt.Sprintf("%x", hash)
}
