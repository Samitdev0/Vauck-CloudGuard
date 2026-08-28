package auth

import (
	"crypto/rand"
	"encoding/base64"

	"golang.org/x/crypto/argon2"
)

func HashPassword(password string) string {

	salt := make([]byte, 16)

	rand.Read(salt)

	hash := argon2.IDKey(
		[]byte(password),
		salt,
		1,
		64*1024,
		4,
		32,
	)

	return base64.RawStdEncoding.EncodeToString(hash)
}
