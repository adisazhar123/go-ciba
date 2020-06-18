package domain

type ClientApplicationInterface interface {
	GetId() string
	SetId(id string)
	GetSecret() string
	SetSecret(id string)
	GetName() string
	SetName(name string)
	GetScope() string
	SetScope(scope string)
	GetTokenMode() string
	SetTokenMode(mode string)
}

type AccessTokenInterface interface {
	IsExpired() bool
	GetScope() string
}

type UserAccountInterface interface {

}


type CibaSessionInterface interface {
	GetClient() ClientApplication
	GetUser() UserAccountInterface
}

type Scope struct {

}

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