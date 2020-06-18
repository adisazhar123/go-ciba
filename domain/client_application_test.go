package domain

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewClientApplication(t *testing.T) {
	ca := NewClientApplication("ca-demo", "openid email profile", "ping", "https://ca-demo.dev/notif", "RS256", false)

	assert.Equal(t, "ca-demo", ca.name)
	assert.Equal(t, "openid email profile", ca.scope)
	assert.Equal(t, "ping", ca.tokenMode)
	assert.Equal(t, "https://ca-demo.dev/notif", ca.clientNotificationEndpoint)
	assert.Equal(t, "RS256", ca.authenticationRequestSigningAlg)
	assert.Equal(t, false, ca.userCodeParameter)
}