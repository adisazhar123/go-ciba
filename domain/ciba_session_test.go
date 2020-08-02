package domain

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewCibaSession(t *testing.T) {
	hint := "some-hint-user-Id"
	bindingMessage := "bind-123"
	token := "someToken-8943dfgdfgdfg5"
	scope := "openid profile Email"
	expiresIn := 120
	interval := 5
	identiferCiba := "urn:openid:params:grant-type:ciba"

	ca := ClientApplication{
		Id:                              "420d637b-ff22-4e48-88fb-237aa2131e72",
		Secret:                          "secret",
		Name:                            "client-app-poll",
		Scope:                           "openid email profile",
		TokenMode:                       ModePoll,
		ClientNotificationEndpoint:      "go-ciba.dev/notification",
		AuthenticationRequestSigningAlg: "",
		UserCodeParameterSupported:      false,
		TokenEndpointAuthMethod:         "client_secret_basic",
		GrantTypes:                      identiferCiba,
	}

	cs := NewCibaSession(&ca, hint, bindingMessage, token, scope, expiresIn, &interval)

	assert.Equal(t, hint, cs.Hint)
	assert.Equal(t, bindingMessage, cs.BindingMessage)
	assert.Equal(t, token, cs.ClientNotificationToken)
	assert.Equal(t, scope, cs.Scope)
	assert.Equal(t, expiresIn, cs.ExpiresIn)
	assert.Equal(t, interval, *cs.Interval)
	assert.NotEmpty(t, cs.AuthReqId)
}
