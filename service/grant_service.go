package service

import "github.com/adisazhar123/go-ciba/util"

type GrantServiceInterface interface {
	ValidateAuthenticationRequestParameters(request *AuthenticationRequest) *util.OidcError
	HandleAuthenticationRequest(request *AuthenticationRequest) (*AuthenticationResponse, *util.OidcError)

	GetGrantIdentifier() string
}

type ConsentServiceInterface interface {
	HandleConsentRequest(request *ConsentRequest)
}
