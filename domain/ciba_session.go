package domain

import (
	"crypto/rand"
	b64 "encoding/base64"
)

type CibaSession struct {
	authReqId string
	expiresIn int
	interval int
}

func NewCibaSession(expiresIn int, interval int) *CibaSession {
	return &CibaSession{
		authReqId: generateAuthReqId(),
		expiresIn: expiresIn,
		interval: interval,
	}
}

func generateAuthReqId() string {
	key := make([]byte, 64)

	_, err := rand.Read(key)
	if err != nil {
		panic("error generating authentication request id")
	}

	return b64.RawURLEncoding.EncodeToString(key)
}