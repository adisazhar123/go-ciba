package http_auth

import (
	"fmt"
	"github.com/adisazhar123/ciba-server/domain"
	"net/http"
)

type ClientAuthenticationContext struct {
	strategy ClientAuthenticationStrategyInterface
}

const ClientSecretBasic = "client_secret_basic"
const ClientSecretPost = "client_secret_post"


func (c *ClientAuthenticationContext) AuthenticateClient(r *http.Request, ca *domain.ClientApplication) bool {
	// TODO: Add more client authentication methods
	switch ca.GetTokenEndpointAuthMethod() {
	case ClientSecretBasic:
		c.strategy = &HttpBasic{clientCredentials: &HttpClientCredentials{}}
	case ClientSecretPost:
		panic(fmt.Sprintf("ciba server doesn't support %s authentication method", ClientSecretPost))
	default:
		panic(fmt.Sprintf("ciba server doesn't support %s authentication method", ca.GetTokenEndpointAuthMethod()))
	}

	return c.AuthenticateClient(r, ca)
}