package http_auth

import (
	"github.com/adisazhar123/go-ciba/domain"
	"net/http"
)

type ClientAuthenticationStrategyInterface interface {
	ValidateRequest(r *http.Request, ca *domain.ClientApplication) bool
}
