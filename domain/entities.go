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
	SetUserCodeSupported(supported bool)
	IsUserCodeSupported() bool
}

type AccessTokenInterface interface {
	IsExpired() bool
	GetScope() string
}

type UserAccountInterface interface {
	GetUseCode() string
	SetUserCode(code string) string
}

type CibaSessionInterface interface {
	GetClient() ClientApplication
	GetUser() UserAccountInterface
}

type Scope struct {

}