package domain

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/adisazhar123/go-ciba/util"
)

type CibaSession struct {
	// This is a unique identifier to identify the authentication request made by the client
	AuthReqId string `db:"auth_req_id"`
	// This is the client application identifier binded to this Ciba session
	ClientId string `db:"client_id"`
	// This is the user identifier the Ciba session is targeting
	UserId string `db:"user_id"`
	// This is the user identifier the Ciba session is targeting
	// as of now only works for static user identifier e.g. user id, email etc.
	Hint string `db:"hint"`
	// This is the binding message/ code to bind session between consumption device/ client application
	// and authentication device.
	BindingMessage string `db:"binding_message"`
	// This is the client notification token that is generated by the client
	// that is used by the Ciba server to communicate as authorization bearer.
	// For token modes ping and push.
	ClientNotificationToken string `db:"client_notification_token"`
	// When the Ciba session (authentication request id) will expire in seconds after the session is created.
	ExpiresIn int64 `db:"expires_in"`
	// The minimum interval rate for poll token requests. Only for token mode poll.
	Interval *int64 `db:"interval"`
	// The validity of the Ciba session.
	Valid bool `db:"valid"`
	// The id token for this Ciba session.
	IdToken string `db:"id_token"`
	// The consent status of this Ciba session (consented/ not consented).
	// User at the authentication device is in charge of the consent.
	Consented *bool `db:"consented"`
	// The scope requested for this Ciba session.
	Scope string `db:"scope"`
	// The latest time a token was requested using this Ciba session.
	// in unix timestamp. Default Value is null, which means it hasn't
	// requested a token yet. This is used for POLL mode only.
	LatestTokenRequestedAt *int64 `db:"latest_token_requested_at"`
	// The time when this Ciba session was created.
	CreatedAt time.Time `db:"created_at"`
}

func (cs *CibaSession) Expire() {
	cs.Valid = false
}

func (cs *CibaSession) IsTimeExpired() bool {
	now := time.Now().UTC()
	seconds, _ := time.ParseDuration(fmt.Sprintf("%ds", cs.ExpiresIn))
	t := cs.CreatedAt.Add(seconds)
	return now.After(t)
}

func (cs *CibaSession) IsConsented() bool {
	return *cs.Consented == true
}

func (cs *CibaSession) IsAuthorizationPending() bool {
	return cs.Consented == nil
}

func (cs *CibaSession) IsValid() bool {
	return cs.Valid == true
}

func generateAuthReqId() string {
	return util.GenerateRandomString()
}

func NewCibaSession(clientApp *ClientApplication, hint, bindingMessage, clientNotificationToken, scope string, expiresIn int64, interval *int64) *CibaSession {
	if clientApp.TokenMode != ModePoll {
		interval = nil
	}
	return &CibaSession{
		Hint:                    hint,
		UserId:                  hint,
		ClientNotificationToken: clientNotificationToken,
		Scope:                   scope,
		BindingMessage:          bindingMessage,
		AuthReqId:               generateAuthReqId(),
		ExpiresIn:               expiresIn,
		Interval:                interval,
		ClientId:                clientApp.Id,
		Valid:                   true,
		CreatedAt:               time.Now().UTC(),
	}
}

func (cs *CibaSession) MarshalBinary() ([]byte, error) {
	return json.Marshal(cs)
}

func (cs *CibaSession) UnmarshalBinary(data []byte) error {
	if err := json.Unmarshal(data, &cs); err != nil {
		return err
	}

	return nil
}
