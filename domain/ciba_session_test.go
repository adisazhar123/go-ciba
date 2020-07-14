package domain

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewCibaSession(t *testing.T) {
	hint := "some-hint-user-id"
	bindingMessage := "bind-123"
	token := "someToken-8943dfgdfgdfg5"
	scope := "openid profile email"
	expiresIn := 120
	interval := 5
	cs := NewCibaSession(hint, bindingMessage, token, scope, expiresIn, interval)

	assert.Equal(t, hint, cs.Hint)
	assert.Equal(t, bindingMessage, cs.BindingMessage)
	assert.Equal(t, token, cs.ClientNotificationToken)
	assert.Equal(t, scope, cs.Scope)
	assert.Equal(t, expiresIn, cs.ExpiresIn)
	assert.Equal(t, interval, cs.Interval)
	assert.NotEmpty(t, cs.AuthReqId)
}
