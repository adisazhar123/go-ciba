package domain

import (
	"encoding/json"
	"github.com/adisazhar123/ciba-server/util"
	"time"
)

type CibaSession struct {
	AuthReqId               string
	ClientId                string
	UserId                  string
	Hint                    string
	BindingMessage          string
	ClientNotificationToken string
	ExpiresIn               int
	Interval                int
	Valid                   bool
	IdToken                 string
	Consented               *bool
	Scope                   string
	LatestTokenRequestedAt  int // in unix timestamp
	CreatedAt               time.Time
}

func generateAuthReqId() string {
	return util.GenerateRandomString()
}

func NewCibaSession(hint, bindingMessage, clientNotificationToken, scope string, expiresIn, interval int) *CibaSession {
	return &CibaSession{
		Hint:                    hint,
		ClientNotificationToken: clientNotificationToken,
		Scope:                   scope,
		BindingMessage:          bindingMessage,
		AuthReqId:               generateAuthReqId(),
		ExpiresIn:               expiresIn,
		Interval:                interval,
		CreatedAt:               time.Now(),
	}
}

func (ca *CibaSession) MarshalBinary() ([]byte, error) {
	return json.Marshal(ca)
}

func (ca *CibaSession) UnmarshalBinary(data []byte) error {
	if err := json.Unmarshal(data, &data); err != nil {
		return err
	}

	return nil
}
