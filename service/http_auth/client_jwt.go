package http_auth

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/adisazhar123/go-ciba/domain"
	"github.com/adisazhar123/go-ciba/util"
)

const jwtBearerAssertionType = "urn:ietf:params:oauth:client-assertion-type:jwt-bearer"

type clientJwt struct {
	goJose                  util.EncryptionInterface
	authServerTokenEndpoint string
}

func (c *clientJwt) GetClientCredentials(r *http.Request, clientId, clientSecret *string) {
	panic("implement me")
}

type claims struct {
	Iss *string    `json:"iss"`
	Sub *string    `json:"sub"`
	Aud *string    `json:"aud"`
	Jti *string    `json:"jti"`
	Exp *time.Time `json:"exp"`
	Iat *time.Time `json:"iat"`
}

func (c *clientJwt) ValidateRequest(r *http.Request, ca *domain.ClientApplication) bool {
	_ = r.ParseForm()
	form := r.Form

	clientAssertion := form.Get("client_assertion")
	if clientAssertion == "" {
		return false
	}
	assertionType := form.Get("client_assertion_type")
	if assertionType != jwtBearerAssertionType {
		return false
	}

	output, err := c.goJose.Decode(clientAssertion, ca.Secret)
	if err != nil {
		return false
	}

	var decodedClaims claims

	err = json.Unmarshal([]byte(output), &decodedClaims)
	if err != nil {
		return false
	}

	if decodedClaims.Iss == nil {
		return false
	}
	if decodedClaims.Sub == nil {
		return false
	}
	if decodedClaims.Aud == nil {
		return false
	}
	if decodedClaims.Jti == nil {
		return false
	}
	if decodedClaims.Exp == nil {
		return false
	}

	return *decodedClaims.Iss == ca.Id && *decodedClaims.Sub == ca.Id && *decodedClaims.Aud == c.authServerTokenEndpoint
}
