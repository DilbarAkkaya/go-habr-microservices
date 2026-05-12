package utils

import (
	"crypto/rand"
	"fmt"
)

func RandomString(length int) string {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return fmt.Sprintf("%d", length)
	}
	return fmt.Sprintf("%x", bytes)[:length]
}
