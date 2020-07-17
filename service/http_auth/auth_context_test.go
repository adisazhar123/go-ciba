package http_auth

import (
	"encoding/base64"
	"github.com/adisazhar123/ciba-server/domain"
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

func TestClientAuthenticationContext_AuthenticateClient_ClientSecretPost_UnupportedType(t *testing.T) {
	method := "POST"
	uri := "ciba.example.com/bc-authorize"
	clientId := "123456"
	clientPassword := "secret"
	auth := base64.StdEncoding.EncodeToString([]byte(clientId + ":" + clientPassword))
	clientApp := &domain.ClientApplication{
		Id:                      clientId,
		Secret:                  clientPassword,
		TokenEndpointAuthMethod: ClientSecretPost,
	}

	req, _ := http.NewRequest(method, uri, nil)
	req.Header.Add("Authorization", "Basic "+auth)

	authContext := &ClientAuthenticationContext{}
	res := authContext.AuthenticateClient(req, clientApp)

	assert.Equal(t, false, res)
}

func TestClientAuthenticationContext_AuthenticateClient_ClientSecretJwt_UnupportedType(t *testing.T) {
	method := "POST"
	uri := "ciba.example.com/bc-authorize"
	clientId := "123456"
	clientPassword := "secret"
	auth := base64.StdEncoding.EncodeToString([]byte(clientId + ":" + clientPassword))
	clientApp := &domain.ClientApplication{
		Id:                      clientId,
		Secret:                  clientPassword,
		TokenEndpointAuthMethod: "client_secret_jwt",
	}

	req, _ := http.NewRequest(method, uri, nil)
	req.Header.Add("Authorization", "Basic "+auth)

	authContext := &ClientAuthenticationContext{}
	res := authContext.AuthenticateClient(req, clientApp)

	assert.Equal(t, false, res)
}
