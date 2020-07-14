package service

type GrantServiceInterface interface {
	ValidateAuthenticationRequestParameters(request *AuthenticationRequest) (interface{}, error)
	HandleAuthenticationRequest(request *AuthenticationRequest) (interface{}, error)

	GetGrantIdentifier() string
}
