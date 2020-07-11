package domain

import (
	"github.com/adisazhar123/ciba-server/util"
	"time"
)

type CibaSession struct {
	authReqId string
	clientId string
	userId string
	hint string
	bindingMessage string
	clientNotificationToken string
	expiresIn int
	interval int
	valid bool
	idToken string
	consented *bool
	scope string
	latestTokenRequestedAt int // in unix timestamp
	createdAt time.Time
}

func generateAuthReqId() string {
	return util.GenerateRandomString()
}

func NewCibaSession(expiresIn int, interval int) *CibaSession {
	return &CibaSession{
		authReqId: generateAuthReqId(),
		expiresIn: expiresIn,
		interval: interval,
		createdAt: time.Now(),
	}
}
