package http_auth

import (
	"net/http"

	"github.com/adisazhar123/go-ciba/domain"
)

type clientPost struct {
}

func (c *clientPost) ValidateRequest(r *http.Request, ca *domain.ClientApplication) bool {
	_ = r.ParseForm()
	form := r.Form

	clientId := form.Get("client_id")
	clientSecret := form.Get("client_secret")

	return ca != nil && ca.Id == clientId && ca.Secret == clientSecret
}

