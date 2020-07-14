package repository

import (
	"github.com/adisazhar123/ciba-server/domain"
	"github.com/alicebob/miniredis"
	"github.com/stretchr/testify/assert"
	"testing"
)

func newTestRedis() *miniredis.Miniredis {
	redis, _ := miniredis.Run()
	return redis
}

func TestClientApplicationRedisRepository_Register(t *testing.T) {
	redis := newTestRedis()

	repo := NewClientApplicationRedisRepository(redis.Addr())
	name := "test-app"
	scope := "openid profile email"
	tokenMode := "ping"
	endpoint := "https://adisazhar.com/notification"
	alg := "RS256"
	userCode := false

	newClientApp := domain.NewClientApplication(name, scope, tokenMode, endpoint, alg, userCode)
	err := repo.Register(newClientApp)

	assert.Empty(t, err)
	assert.Equal(t, newClientApp.GetName(), name)
	assert.Equal(t, newClientApp.GetScope(), scope)
	assert.Equal(t, newClientApp.GetTokenMode(), tokenMode)
	assert.Equal(t, newClientApp.GetClientNotificationEndpoint(), endpoint)
	assert.Equal(t, newClientApp.GetAuthenticationRequestSigningAlg(), alg)
	assert.Equal(t, newClientApp.GetUserCodeParameterSupported(), userCode)
}

func TestCibaSessionRedisRepository_Create(t *testing.T) {
	redis := newTestRedis()
	repo := NewCibaSessionRedisRepository(redis.Addr())
	hint := "some-hint-user-id"
	bindingMessage := "bind-123"
	token := "someToken-8943dfgdfgdfg5"
	scope := "openid profile email"
	expiresIn := 120
	interval := 5

	newCibaSession := domain.NewCibaSession(hint, bindingMessage, token, scope, expiresIn, interval)
	err := repo.Create(newCibaSession)

	assert.Empty(t, err)
	assert.Equal(t, newCibaSession.Hint, hint)
	assert.Equal(t, newCibaSession.BindingMessage, bindingMessage)
	assert.Equal(t, newCibaSession.ClientNotificationToken, token)
	assert.Equal(t, newCibaSession.Scope, scope)
	assert.Equal(t, newCibaSession.ExpiresIn, expiresIn)
	assert.Equal(t, newCibaSession.Interval, interval)
}
