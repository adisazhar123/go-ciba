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