package repository

import (
	"fmt"
	"testing"
	"time"

	"github.com/adisazhar123/go-ciba/domain"
	"github.com/adisazhar123/go-ciba/grant"
	"github.com/adisazhar123/go-ciba/test_data"
	"github.com/alicebob/miniredis"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
)

func newTestRedis() *miniredis.Miniredis {
	miniRedis, _ := miniredis.Run()
	return miniRedis
}

func newRedisClient(addr string) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:               addr,
	})
}

func TestClientApplicationRedisRepository_Register(t *testing.T) {
	miniRedis := newTestRedis()

	repo := NewClientApplicationRedisRepository(newRedisClient(miniRedis.Addr()))
	name := "test-app"
	scope := "openid profile email"
	tokenMode := "ping"
	endpoint := "https://adisazhar.com/notification"
	alg := "RS256"
	userCode := false

	newClientApp := domain.NewClientApplication(name, scope, tokenMode, endpoint, alg, userCode)
	marshalled, _ := newClientApp.MarshalBinary()

	err := repo.Register(newClientApp)

	miniRedis.CheckGet(t, "client_application:"+newClientApp.Id, string(marshalled))
	assert.NoError(t, err)
}

func TestClientApplicationRedisRepository_FindById(t *testing.T) {
	miniRedis := newTestRedis()
	repo := NewClientApplicationRedisRepository(newRedisClient(miniRedis.Addr()))
	name := "test-app"
	scope := "openid profile email"
	tokenMode := "ping"
	endpoint := "https://adisazhar.com/notification"
	alg := "RS256"
	userCode := false
	newClientApp := domain.NewClientApplication(name, scope, tokenMode, endpoint, alg, userCode)
	jsonString, _ := newClientApp.MarshalBinary()
	miniRedis.Set("client_application:"+newClientApp.Id, string(jsonString))

	clientApp, err := repo.FindById(newClientApp.Id)

	assert.NotNil(t, clientApp)
	assert.NoError(t, err)
}

func TestCibaSessionRedisRepository_Create(t *testing.T) {
	miniRedis := newTestRedis()
	repo := NewCibaSessionRedisRepository(newRedisClient(miniRedis.Addr()))
	hint := "some-hint-user-id"
	bindingMessage := "bind-123"
	token := "someToken-8943dfgdfgdfg5"
	scope := "openid profile email"
	expiresIn := int64(120)
	interval := int64(5)

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
	marshalled, _ := newCibaSession.MarshalBinary()

	err := repo.Create(newCibaSession)

	miniRedis.CheckGet(t, "ciba_session:"+newCibaSession.AuthReqId, string(marshalled))
	assert.NoError(t, err)
}

func TestUserAccountRedisRepository_FindById_ValidUser(t *testing.T) {
	miniRedis := newTestRedis()
	repo := NewUserAccountRedisRepository(newRedisClient(miniRedis.Addr()))
	userId := "23246440-92d9-4738-8faf-551d24a1c4a4"
	user := &domain.UserAccount{
		Id:        "23246440-92d9-4738-8faf-551d24a1c4a4",
		Name:      "User Account 01",
		Email:     "ua@email.com",
		Password:  "secret",
		UserCode:  "123",
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
	jsonString, _ := user.MarshalBinary()
	miniRedis.Set("user_account:"+userId, string(jsonString))

	foundUser, err := repo.FindById(userId)

	assert.Empty(t, err)
	assert.Equal(t, user.Id, foundUser.Id)
}

func TestUserAccountRedisRepository_FindById_InvalidUser(t *testing.T) {
	miniRedis := newTestRedis()
	repo := NewUserAccountRedisRepository(newRedisClient(miniRedis.Addr()))
	invalidUserId := "invalid"
	user := &domain.UserAccount{
		Id:        "23246440-92d9-4738-8faf-551d24a1c4a4",
		Name:      "User Account 01",
		Email:     "ua@email.com",
		Password:  "secret",
		UserCode:  "123",
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
	jsonString, _ := user.MarshalBinary()
	miniRedis.Set("user_account:"+user.Id, string(jsonString))

	foundUser, err := repo.FindById(invalidUserId)

	assert.Nil(t, foundUser)
	assert.NoError(t, err)
}

func TestCibaSessionRedisRepository_FindById_ShouldReturnNil_WhenNotFound(t *testing.T) {
	miniRedis := newTestRedis()
	repo := NewCibaSessionRedisRepository(newRedisClient(miniRedis.Addr()))
	invalidSessionId := "invalid"
	bytes, _ := test_data.CibaSession1.MarshalBinary()
	miniRedis.Set("ciba_session:"+test_data.CibaSession1.AuthReqId, string(bytes))

	ciba, err := repo.FindById(invalidSessionId)

	assert.Nil(t, ciba)
	assert.NoError(t, err)
}

func TestCibaSessionRedisRepository_FindById_ShouldReturnCibaSession(t *testing.T) {
	miniRedis := newTestRedis()
	repo := NewCibaSessionRedisRepository(newRedisClient(miniRedis.Addr()))
	bytes, _ := test_data.CibaSession1.MarshalBinary()
	miniRedis.Set("ciba_session:"+test_data.CibaSession1.AuthReqId, string(bytes))

	cs, err := repo.FindById(test_data.CibaSession1.AuthReqId)

	assert.Nil(t, err)
	assert.NotNil(t, cs)
}

func TestCibaSessionRedisRepository_Update(t *testing.T) {
	miniRedis := newTestRedis()
	repo := NewCibaSessionRedisRepository(newRedisClient(miniRedis.Addr()))
	cs := test_data.CibaSession1
	bytes, _ := test_data.CibaSession1.MarshalBinary()

	err := repo.Update(&cs)

	miniRedis.CheckGet(t, "ciba_session:"+test_data.CibaSession1.AuthReqId, string(bytes))
	assert.Nil(t, err)
}

func TestKeyRedisRepository_FindPrivateKeyByClientId(t *testing.T) {
	miniRedis := newTestRedis()
	repo := NewKeyRedisRepository(newRedisClient(miniRedis.Addr()))
	bytes, _ := test_data.Key1.MarshalBinary()
	miniRedis.Set("oauth_key:"+test_data.Key1.ClientId, string(bytes))

	key, err := repo.FindPrivateKeyByClientId(test_data.Key1.ClientId)

	assert.Nil(t, err)
	assert.Equal(t, test_data.Key1, *key)
}

func TestAccessTokenRedisRepository_Create(t *testing.T) {
	miniRedis := newTestRedis()
	repo := NewAccessTokenRedisRepository(newRedisClient(miniRedis.Addr()))
	accessToken := domain.NewAccessToken("1-1-1-1", "2-2-2-2", "3-3-3-3", "openid address", time.Now().UTC().Add(1 * time.Hour))
	marshalled, _ := accessToken.MarshalBinary()

	err := repo.Create(accessToken)

	miniRedis.CheckGet(t, "access_token:"+accessToken.Value, string(marshalled))
	assert.NoError(t, err)
}

func TestAccessTokenRedisRepository_Find(t *testing.T) {
	miniRedis := newTestRedis()
	repo := NewAccessTokenRedisRepository(newRedisClient(miniRedis.Addr()))
	accessToken := domain.NewAccessToken("1-1-1-1", "2-2-2-2", "3-3-3-3", "openid address", time.Now().UTC().Add(1 * time.Hour))
	marshalled, _ := accessToken.MarshalBinary()
	miniRedis.Set("access_token:"+accessToken.Value, string(marshalled))

	at, err := repo.Find(accessToken.Value)

	assert.NotNil(t, at)
	assert.NoError(t, err)
}

func TestUserClaimRedisRepository_GetUserClaims(t *testing.T) {
	miniRedis := newTestRedis()
	repo := NewUserClaimRedisRepository(newRedisClient(miniRedis.Addr()))
	userAccount := test_data.User1
	marshalled, _ := userAccount.MarshalBinary()
	miniRedis.Set("user_account:"+userAccount.Id, string(marshalled))
	miniRedis.Lpush("scope:openid", "id")


	claims, err := repo.GetUserClaims(userAccount.Id, "openid address")

	miniRedis.CheckList(t, "scope:openid", "id")
	miniRedis.CheckGet(t, "user_account:"+userAccount.Id, string(marshalled))
	assert.NotNil(t, claims)
	assert.NoError(t, err)
	assert.Contains(t, claims, "id")
}