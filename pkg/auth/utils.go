package auth

import (
	"crypto/sha1"
	"fmt"
)

// generatePasswordHash генерирует hash пароля.
func GeneratePasswordHash(password string, salt string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}
