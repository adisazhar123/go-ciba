package http_auth

import (
	"net/http"

	"github.com/adisazhar123/go-ciba/domain"
	"github.com/adisazhar123/go-ciba/grant"
	"github.com/adisazhar123/go-ciba/util"
)

type ClientAuthenticationContext struct {
	strategy    ClientAuthenticationStrategyInterface
	grantConfig *grant.GrantConfig
}

func NewClientAuthenticationContext(grantConfig *grant.GrantConfig) *ClientAuthenticationContext {
	return &ClientAuthenticationContext{
		grantConfig: grantConfig,
	}
}

const (
	ClientSecretBasic = "client_secret_basic"
	ClientSecretPost  = "client_secret_post"
	ClientSecretJwt   = "client_secret_jwt"
)

var supportedClientAuthentications map[string]ClientAuthenticationStrategyInterface = map[string]ClientAuthenticationStrategyInterface{
	ClientSecretBasic: &httpBasic{clientCredentials: &httpClientCredentials{}},
	ClientSecretPost:  &clientPost{},
}

func PopulateClientCredentials(r *http.Request, clientId, clientSecret *string)  {
	for _, v := range supportedClientAuthentications {
		v.GetClientCredentials(r, clientId, clientSecret)
		if *clientId != "" && *clientSecret != "" {
			break
		}
	}
}

func (c *ClientAuthenticationContext) AuthenticateClient(r *http.Request, ca *domain.ClientApplication) bool {
	// If no method is registered, the default method is client_secret_basic
	switch ca.GetTokenEndpointAuthMethod() {
	case ClientSecretBasic:
		c.strategy = &httpBasic{clientCredentials: &httpClientCredentials{}}
	case ClientSecretPost:
		c.strategy = &clientPost{}
	case ClientSecretJwt:
		c.strategy = &clientJwt{
			goJose:                  util.NewGoJoseEncryption(),
			authServerTokenEndpoint: c.grantConfig.TokenEndpointUrl,
		}
	default:
		c.strategy = &httpBasic{clientCredentials: &httpClientCredentials{}}
	}

	return c.strategy.ValidateRequest(r, ca)
}
