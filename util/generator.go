package util

import (
	"crypto/rand"
	b64 "encoding/base64"
	"github.com/google/uuid"
)

func GenerateRandomString() string {
	key := make([]byte, 64)

	_, err := rand.Read(key)
	if err != nil {
		panic("error generating authentication request id")
	}

	return b64.RawURLEncoding.EncodeToString(key)
}

func GenerateUuid() string {
	return uuid.New().String()
}