package domain

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewClientApplication(t *testing.T) {
	ca := NewClientApplication("ca-demo", "openid email profile", "ping", "https://ca-demo.dev/notif", "RS256", false)

	assert.Equal(t, "ca-demo", ca.Name)
	assert.Equal(t, "openid email profile", ca.Scope)
	assert.Equal(t, "ping", ca.TokenMode)
	assert.Equal(t, "https://ca-demo.dev/notif", ca.ClientNotificationEndpoint)
	assert.Equal(t, "RS256", ca.AuthenticationRequestSigningAlg)
	assert.Equal(t, false, ca.UserCodeParameter)
}