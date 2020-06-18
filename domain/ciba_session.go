package domain

import (
	"github.com/adisazhar123/ciba-server/util"
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
	return util.GenerateRandomString()
}
