package http_auth

import (
	"encoding/base64"
	"github.com/adisazhar123/ciba-server/domain"
	"log"
	"net/http"
	"strings"
)

type HttpBasic struct {
	clientCredentials *HttpClientCredentials
}

type HttpClientCredentials struct {
	ClientId     string
	ClientSecret string
}

func (hb *HttpBasic) ValidateRequest(r *http.Request, ca *domain.ClientApplication) bool {
	clientCred := hb.getClientCredentials(r)
	log.Println(ca)
	return clientCred != nil && ca.GetId() == clientCred.ClientId && ca.GetSecret() == clientCred.ClientSecret
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
		ClientId:     credentials[0],
		ClientSecret: credentials[1],
	}
}

func UtilGetClientCredentials(r *http.Request) *HttpClientCredentials {
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
		ClientId:     credentials[0],
		ClientSecret: credentials[1],
	}
}
