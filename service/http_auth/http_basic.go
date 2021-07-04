package http_auth

import (
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/adisazhar123/go-ciba/domain"
)

type httpBasic struct {
	clientCredentials *httpClientCredentials
}

func (hb *httpBasic) GetClientCredentials(r *http.Request, clientId, clientSecret *string) {
	credentials := hb.getClientCredentials(r)
	if credentials != nil {
		*clientId = credentials.clientId
		*clientSecret = credentials.clientSecret
	}
}

func (hb *httpBasic) ValidateRequest(r *http.Request, ca *domain.ClientApplication) bool {
	clientCred := hb.getClientCredentials(r)
	return clientCred != nil && ca != nil && ca.GetId() == clientCred.clientId && ca.GetSecret() == clientCred.clientSecret
}

func (hb *httpBasic) getClientCredentials(r *http.Request) *httpClientCredentials {
	return UtilGetClientCredentials(r)
}

func UtilGetClientCredentials(r *http.Request) *httpClientCredentials {
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

	return &httpClientCredentials{
		clientId:     credentials[0],
		clientSecret: credentials[1],
	}
}
