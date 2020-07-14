package http_auth

import (
	"github.com/adisazhar123/ciba-server/domain"
	"net/http"
)

type ClientAuthenticationStrategyInterface interface {
	ValidateRequest(r *http.Request, ca *domain.ClientApplication) bool
}
