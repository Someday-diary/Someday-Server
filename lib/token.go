package lib

import (
	"crypto/rand"
	"fmt"
)

func CreateToken(size int) string {
	b := make([]byte, size)
	_, _ = rand.Read(b)
	key := fmt.Sprintf("%x", b)[:size]
	return key
}
