package service

import "github.com/adisazhar123/go-ciba/util"

type GrantServiceInterface interface {
	ValidateAuthenticationRequestParameters(request *AuthenticationRequest) (interface{}, *util.OidcError)
	HandleAuthenticationRequest(request *AuthenticationRequest) (interface{}, *util.OidcError)

	GetGrantIdentifier() string
}

type ConsentServiceInterface interface {
	HandleConsentRequest(request *ConsentRequest)
}
