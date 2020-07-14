package util

import (
	"crypto/rand"
	b64 "encoding/base64"
	"github.com/google/uuid"
)

func GenerateRandomString() string {
	key := make([]byte, 64)
	_, _ = rand.Read(key)

	return b64.RawURLEncoding.EncodeToString(key)
}

func GenerateUuid() string {
	return uuid.New().String()
}
