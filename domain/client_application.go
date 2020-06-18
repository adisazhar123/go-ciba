package domain

import "github.com/adisazhar123/ciba-server/util"

type ClientApplication struct {
	id string
	secret string
	name string
	scope string
	tokenMode string
	clientNotificationEndpoint string
	authenticationRequestSigningAlg string
	userCodeParameter bool
}

func NewClientApplication(name, scope, tokenMode, clientNotificationEndpoint, authenticationRequestSigningAlg string, userCode bool) *ClientApplication {
	return &ClientApplication{
		id:                              util.GenerateUuid(),
		secret:                          util.GenerateRandomString(),
		name:                            name,
		scope:                           scope,
		tokenMode:                       tokenMode,
		clientNotificationEndpoint:      clientNotificationEndpoint,
		authenticationRequestSigningAlg: authenticationRequestSigningAlg,
		userCodeParameter:               userCode,
	}
}

func (ca *ClientApplication) SetId(id string) {
	ca.id = id
}

func (ca *ClientApplication) SetSecret(secret string) {
	ca.secret = secret
}

func (ca *ClientApplication) SetName(name string) {
	ca.name = name
}

func (ca *ClientApplication) SetScope(scope string) {
	ca.scope = scope
}

func (ca *ClientApplication) SetTokenMode(mode string) {
	ca.tokenMode = mode
}