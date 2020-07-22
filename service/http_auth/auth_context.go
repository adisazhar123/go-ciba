package http_auth

import (
	"fmt"
	"github.com/adisazhar123/go-ciba/domain"
	"log"
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
		log.Println(fmt.Sprintf("ciba server doesn't support %s authentication method", ca.GetTokenEndpointAuthMethod()))
		return false
	default:
		log.Println(fmt.Sprintf("ciba server doesn't support %s authentication method", ca.GetTokenEndpointAuthMethod()))
		return false
	}

	return c.strategy.ValidateRequest(r, ca)
}
