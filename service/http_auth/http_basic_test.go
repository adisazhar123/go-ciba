package http_auth

import (
	"encoding/base64"
	"github.com/adisazhar123/ciba-server/domain"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestHttpBasic_ValidateRequest_CorrectClientCredentials(t *testing.T) {
	httpBasic := &HttpBasic{}
	method := "POST"
	uri := "ciba.example.com/bc-authorize"

	clientApp := &domain.ClientApplication{
		Id: "123456",
		Secret: "secret",
	}
	auth := base64.StdEncoding.EncodeToString([]byte(clientApp.Id + ":" + clientApp.Secret))
	req, _ := http.NewRequest(method, uri, nil)
	req.Header.Add("Authorization", "Basic " + auth)

	res := httpBasic.ValidateRequest(req, clientApp)

	assert.Equal(t, true, res)
}

func TestHttpBasic_ValidateRequest_IncorrectClientCredentials(t *testing.T) {
	httpBasic := &HttpBasic{}
	method := "POST"
	uri := "ciba.example.com/bc-authorize"

	clientApp := &domain.ClientApplication{
		Id: "123456",
		Secret: "secret",
	}
	auth := base64.StdEncoding.EncodeToString([]byte(clientApp.Id + ":" + clientApp.Secret + "extra-password"))
	req, _ := http.NewRequest(method, uri, nil)
	req.Header.Add("Authorization", "Basic " + auth)

	res := httpBasic.ValidateRequest(req, clientApp)

	assert.Equal(t, false, res)
}

func TestHttpBasic_GetClientCredentials_Valid(t *testing.T) {
	httpBasic := &HttpBasic{}
	method := "POST"
	uri := "ciba.example.com/bc-authorize"
	clientId := "123456"
	clientPassword := "secret"
	auth := base64.StdEncoding.EncodeToString([]byte(clientId + ":" + clientPassword))

	req, _ := http.NewRequest(method, uri, nil)
	req.Header.Add("Authorization", "Basic " + auth)

	credentials := httpBasic.getClientCredentials(req)

	assert.Equal(t, clientId, credentials.clientId)
	assert.Equal(t, clientPassword, credentials.clientSecret)
}

func TestHttpBasic_GetClientCredentials_AuthorizationValueIncorrectlyFormed(t *testing.T) {
	httpBasic := &HttpBasic{}
	method := "POST"
	uri := "ciba.example.com/bc-authorize"
	clientId := "123456"
	clientPassword := "secret"
	auth := base64.StdEncoding.EncodeToString([]byte(clientId + clientPassword))

	req, _ := http.NewRequest(method, uri, nil)
	req.Header.Add("Authorization", "Basic " + auth)

	credentials := httpBasic.getClientCredentials(req)

	assert.Nil(t, credentials)
}

func TestHttpBasic_GetClientCredentials_AuthorizationValueIncorrectlyFormed2(t *testing.T) {
	httpBasic := &HttpBasic{}
	method := "POST"
	uri := "ciba.example.com/bc-authorize"
	clientId := "123456"
	clientPassword := "secret"
	auth := base64.StdEncoding.EncodeToString([]byte(clientId + ":" + clientPassword + ":" + "extra"))

	req, _ := http.NewRequest(method, uri, nil)
	req.Header.Add("Authorization", "Basic " + auth)

	credentials := httpBasic.getClientCredentials(req)

	assert.Nil(t, credentials)
}

func TestHttpBasic_GetClientCredentials_AuthorizationValueIncorrectlyFormed3(t *testing.T) {
	httpBasic := &HttpBasic{}
	method := "POST"
	uri := "ciba.example.com/bc-authorize"
	clientId := "123456"
	clientPassword := "secret"
	auth := base64.StdEncoding.EncodeToString([]byte(clientId + ":" + clientPassword))

	req, _ := http.NewRequest(method, uri, nil)
	req.Header.Add("Authorization", auth)

	credentials := httpBasic.getClientCredentials(req)

	assert.Nil(t, credentials)
}

func TestHttpBasic_GetClientCredentials_AuthorizationValueIncorrectEncoding(t *testing.T) {
	httpBasic := &HttpBasic{}
	method := "POST"
	uri := "ciba.example.com/bc-authorize"
	clientId := "123456"
	clientPassword := "secret"
	auth := base64.StdEncoding.EncodeToString([]byte(clientId + ":" + clientPassword))

	req, _ := http.NewRequest(method, uri, nil)
	req.Header.Add("Authorization", "Basic " + auth + "make-it-incorrect")

	credentials := httpBasic.getClientCredentials(req)

	assert.Nil(t, credentials)
}