package grant

type GrantTypeInterface interface {
	GetIdentifier() string
	ValidateAuthenticationRequest()
	HandleAuthenticationRequest()
}