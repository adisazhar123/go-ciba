package http_auth

import (
	"encoding/base64"
	"github.com/adisazhar123/ciba-server/domain"
	"net/http"
	"strings"
)

type HttpBasic struct {
	clientCredentials *HttpClientCredentials
}

type HttpClientCredentials struct {
	clientId     string
	clientSecret string
}

func (hb *HttpBasic) ValidateRequest(r *http.Request, ca *domain.ClientApplication) bool {
	clientCred := hb.getClientCredentials(r)
	return clientCred != nil && ca.GetId() == clientCred.clientId && ca.GetSecret() == clientCred.clientSecret
}

func (hb *HttpBasic) getClientCredentials(r *http.Request) *HttpClientCredentials {
	header := r.Header.Get("Authorization")
	splitToken := strings.Split(header, "Basic")
	if len(splitToken) != 2 {
		// TODO: Add logging http header not correct
		return nil
	}
	encodedToken := strings.TrimSpace(splitToken[1])
	decodedToken, err := base64.StdEncoding.DecodeString(encodedToken)
	if err != nil {
		// TODO: Add logging http header token not correct
		return nil
	}

	credentials := strings.Split(string(decodedToken), ":")
	if len(credentials) != 2 {
		// TODO: Add logging token not correct
		return nil
	}

	return &HttpClientCredentials{
		clientId:     credentials[0],
		clientSecret: credentials[1],
	}
}
