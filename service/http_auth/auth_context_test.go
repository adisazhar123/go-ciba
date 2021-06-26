package http_auth

import (
	"encoding/base64"
	"github.com/adisazhar123/go-ciba/domain"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestClientAuthenticationContext_AuthenticateClient_ClientSecretBasic_SupportedType(t *testing.T) {
	method := "POST"
	uri := "ciba.example.com/bc-authorize"
	clientId := "123456"
	clientPassword := "secret"
	auth := base64.StdEncoding.EncodeToString([]byte(clientId + ":" + clientPassword))
	clientApp := &domain.ClientApplication{
		Id:                      clientId,
		Secret:                  clientPassword,
		TokenEndpointAuthMethod: ClientSecretBasic,
	}

	req, _ := http.NewRequest(method, uri, nil)
	req.Header.Add("Authorization", "Basic "+auth)

	authContext := &ClientAuthenticationContext{}
	res := authContext.AuthenticateClient(req, clientApp)

	assert.Equal(t, true, res)
}
