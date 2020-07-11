package grant

import "github.com/adisazhar123/ciba-server/service"

type GrantTypeInterface interface {
	GetIdentifier() string
	ValidateAuthenticationRequest(request *service.AuthenticationRequest) error
	HandleAuthenticationRequest(request *service.AuthenticationRequest)
}