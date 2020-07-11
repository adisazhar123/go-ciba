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
	GetId() string
	SetId(id string)
	SetName(name string)
	GetName() string
	SetEmail(email string)
	GetEmail() string
	SetPassword(password string)
	GetPassword() string
	GetUseCode() string
	SetUserCode(code string)
}

type CibaSessionInterface interface {
	GetClient() ClientApplication
	GetUser() UserAccount
}



type Scope struct {

}