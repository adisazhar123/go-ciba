package repository

import (
	"fmt"
	"github.com/adisazhar123/go-ciba/domain"
	"github.com/adisazhar123/go-ciba/grant"
	"github.com/alicebob/miniredis"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
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

	ca := domain.ClientApplication{
		Id:                              "420d637b-ff22-4e48-88fb-237aa2131e72",
		Secret:                          "secret",
		Name:                            "client-app-poll",
		Scope:                           "openid email profile",
		TokenMode:                       domain.ModePoll,
		ClientNotificationEndpoint:      "go-ciba.dev/notification",
		AuthenticationRequestSigningAlg: "",
		UserCodeParameterSupported:      false,
		TokenEndpointAuthMethod:         "client_secret_basic",
		GrantTypes:                      fmt.Sprintf("%s", grant.IdentifierCiba),
	}

	newCibaSession := domain.NewCibaSession(&ca, hint, bindingMessage, token, scope, expiresIn, &interval)
	err := repo.Create(newCibaSession)

	assert.Empty(t, err)
	assert.Equal(t, newCibaSession.Hint, hint)
	assert.Equal(t, newCibaSession.BindingMessage, bindingMessage)
	assert.Equal(t, newCibaSession.ClientNotificationToken, token)
	assert.Equal(t, newCibaSession.Scope, scope)
	assert.Equal(t, newCibaSession.ExpiresIn, expiresIn)
	assert.Equal(t, *newCibaSession.Interval, interval)
}

func TestUserAccountRedisRepository_FindById_ValidUser(t *testing.T) {
	redis := newTestRedis()
	repo := NewUserAccountRedisRepository(redis.Addr())
	userId := "23246440-92d9-4738-8faf-551d24a1c4a4"
	user := &domain.UserAccount{
		Id:        "23246440-92d9-4738-8faf-551d24a1c4a4",
		Name:      "User Account 01",
		Email:     "ua@email.com",
		Password:  "secret",
		UserCode:  "123",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	jsonString, _ := user.MarshalBinary()
	redis.Set("user_account:"+userId, string(jsonString))

	foundUser, err := repo.FindById(userId)

	assert.Empty(t, err)
	assert.Equal(t, user.Id, foundUser.Id)
}

func TestUserAccountRedisRepository_FindById_InvalidUser(t *testing.T) {
	redis := newTestRedis()
	repo := NewUserAccountRedisRepository(redis.Addr())
	invalidUserId := "invalid"
	user := &domain.UserAccount{
		Id:        "23246440-92d9-4738-8faf-551d24a1c4a4",
		Name:      "User Account 01",
		Email:     "ua@email.com",
		Password:  "secret",
		UserCode:  "123",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	jsonString, _ := user.MarshalBinary()
	redis.Set("user_account:"+user.Id, string(jsonString))

	foundUser, err := repo.FindById(invalidUserId)

	assert.Empty(t, foundUser)
	assert.NotNil(t, err)
}
