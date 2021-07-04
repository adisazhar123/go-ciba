package http_auth

import (
	"github.com/adisazhar123/go-ciba/domain"
	"net/http"
)

type ClientAuthenticationStrategyInterface interface {
	ValidateRequest(r *http.Request, ca *domain.ClientApplication) bool
}

type httpClientCredentials struct {
	clientId     string
	clientSecret string
}

func (h *httpClientCredentials) GetClientId() string {
	return h.clientId
}

func (h *httpClientCredentials) GetClientSecret() string {
	return h.clientSecret
}
