package domain

import (
	"encoding/json"
	"github.com/adisazhar123/ciba-server/util"
)


const (
	MODE_PING = "ping"
	MODE_POLL = "poll"
	MODE_PUSH = "push"
)

type ClientApplication struct {
	Id                              string
	Secret                          string
	Name                            string
	Scope                           string
	TokenMode                       string
	ClientNotificationEndpoint      string
	AuthenticationRequestSigningAlg string
	UserCodeParameter               bool
}

func NewClientApplication(name, scope, tokenMode, clientNotificationEndpoint, authenticationRequestSigningAlg string, userCode bool) *ClientApplication {
	return &ClientApplication{
		Id:                              util.GenerateUuid(),
		Secret:                          util.GenerateRandomString(),
		Name:                            name,
		Scope:                           scope,
		TokenMode:                       tokenMode,
		ClientNotificationEndpoint:      clientNotificationEndpoint,
		AuthenticationRequestSigningAlg: authenticationRequestSigningAlg,
		UserCodeParameter:               userCode,
	}
}

func (ca *ClientApplication) MarshalBinary() ([]byte, error) {
	return json.Marshal(ca)
}

func (ca *ClientApplication) UnmarshalBinary(data []byte) error {
	if err := json.Unmarshal(data, &data); err != nil {
		return err
	}

	return nil
}

func (ca *ClientApplication) GetId() string {
	return ca.Id
}

func (ca *ClientApplication) GetSecret() string {
	return ca.Secret
}

func (ca *ClientApplication) GetName() string {
	return ca.Name
}

func (ca *ClientApplication) GetScope() string {
	return ca.Scope
}

func (ca *ClientApplication) GetTokenMode() string {
	return ca.TokenMode
}

func (ca *ClientApplication) SetId(id string) {
	ca.Id = id
}

func (ca *ClientApplication) SetSecret(secret string) {
	ca.Secret = secret
}

func (ca *ClientApplication) SetName(name string) {
	ca.Name = name
}

func (ca *ClientApplication) SetScope(scope string) {
	ca.Scope = scope
}

func (ca *ClientApplication) SetTokenMode(mode string) {
	ca.TokenMode = mode
}

func (ca *ClientApplication) SetUserCodeSupported(supported bool) {
	ca.UserCodeParameter = supported
}

func (ca *ClientApplication) IsUserCodeSupported() bool {
	return ca.UserCodeParameter
}