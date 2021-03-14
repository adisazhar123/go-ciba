package service

import (
	"testing"

	"github.com/adisazhar123/go-ciba/domain"
	"github.com/adisazhar123/go-ciba/grant"
	"github.com/adisazhar123/go-ciba/test_data"
	"github.com/adisazhar123/go-ciba/util"
	"github.com/stretchr/testify/assert"
)

type AccessTokenVolatileRepository struct {
	data map[string]*domain.AccessToken
}

func newAccessTokenVolatileRepository() *AccessTokenVolatileRepository {
	return &AccessTokenVolatileRepository{data: map[string]*domain.AccessToken{}}
}

func (a *AccessTokenVolatileRepository) Create(accessToken *domain.AccessToken) error {
	a.data[accessToken.Value] = accessToken
	return nil
}

func newTokenService() *TokenService {
	return &TokenService{
		accessTokenRepo: newAccessTokenVolatileRepository(),
		clientAppRepo:   test_data.NewClientApplicationVolatileRepository(),
		cibaSessionRepo: test_data.NewCibaSessionVolatileRepository(),
		keyRepo:         test_data.NewKeyVolatileRepository(),
		config: &TokenConfig{
			PollingInterval:     5,
			AccessTokenLifeTime: 3600,
			IdTokenLifeTime:     3600,
			Issuer:              "issuer-ciba.example.com",
		},
		grant: grant.NewCibaGrant(),
	}
}

func TestTokenService_GrantAccessToken_ShouldReturnErrorInvalidGrant_WhenAuthReqIdDoesntExist(t *testing.T) {
	ts := newTokenService()

	_, err := ts.GrantAccessToken(&TokenRequest{
		authReqId: "unknown-auth-req-id",
	})

	assert.EqualError(t, err, util.ErrInvalidGrant.Error())
}

func TestTokenService_GrantAccessToken_ShouldReturnErrorInvalidGrant_WhenAuthReqIdIsIssuedToAnotherClient(t *testing.T) {
	ts := newTokenService()

	_, err := ts.GrantAccessToken(&TokenRequest{
		clientId:  "different-client-id",
		authReqId: test_data.CibaSession1.AuthReqId,
	})

	assert.EqualError(t, err, util.ErrInvalidGrant.Error())
}

func TestTokenService_GrantAccessToken_ShouldReturnErrorInvalidClient_WhenAuthReqIdIsntAttachedToClientApplication(t *testing.T) {
	ts := newTokenService()

	_, err := ts.GrantAccessToken(&TokenRequest{
		clientId:  test_data.CibaSession2.ClientId,
		authReqId: test_data.CibaSession2.AuthReqId,
	})

	assert.EqualError(t, err, util.ErrInvalidClient.Error())
}

func TestTokenService_GrantAccessToken_ShouldReturnErrorUnauthorizedClient_WhenClientIsntRegisteredToUseCiba(t *testing.T) {
	ts := newTokenService()

	_, err := ts.GrantAccessToken(&TokenRequest{
		clientId:  test_data.ClientAppNotRegisteredToUseCiba.Id,
		authReqId: test_data.CibaSession3.AuthReqId,
	})

	assert.EqualError(t, err, util.ErrUnauthorizedClient.Error())
}

func TestTokenService_GrantAccessToken_ShouldReturnErrorUnauthorizedClient_WhenClientIsRegisteredToUsePushTokenMode(t *testing.T) {
	ts := newTokenService()

	_, err := ts.GrantAccessToken(&TokenRequest{
		clientId:  test_data.CibaSession4.ClientId,
		authReqId: test_data.CibaSession4.AuthReqId,
	})

	assert.EqualError(t, err, util.ErrUnauthorizedClient.Error())
}

func TestTokenService_GrantAccessToken_ShouldReturnErrorExpiredToken_WhenTokenIsExpired(t *testing.T) {
	ts := newTokenService()

	_, err := ts.GrantAccessToken(&TokenRequest{
		clientId:  test_data.CibaSession5.ClientId,
		authReqId: test_data.CibaSession5.AuthReqId,
	})

	assert.EqualError(t, err, util.ErrExpiredToken.Error())
}

func TestTokenService_GrantAccessToken_ShouldReturnErrorAuthorizationPending_WhenUserHasntGivenConsent(t *testing.T) {
	ts := newTokenService()

	_, err := ts.GrantAccessToken(&TokenRequest{
		clientId:  test_data.CibaSession6.ClientId,
		authReqId: test_data.CibaSession6.AuthReqId,
	})

	assert.EqualError(t, err, util.ErrAuthorizationPending.Error())
}

func TestTokenService_GrantAccessToken_ShouldReturnErrorAccessDenied_WhenUserDidntGiveConsent(t *testing.T) {
	ts := newTokenService()

	_, err := ts.GrantAccessToken(&TokenRequest{
		clientId:  test_data.CibaSession7.ClientId,
		authReqId: test_data.CibaSession7.AuthReqId,
	})

	assert.EqualError(t, err, util.ErrAccessDenied.Error())
}

func TestTokenService_GrantAccessToken_ShouldReturnError_WhenKeyIsntFoundForClientApp(t *testing.T) {
	ts := newTokenService()
	req := &TokenRequest{
		clientId:  test_data.CibaSession8.ClientId,
		authReqId: test_data.CibaSession8.AuthReqId,
	}
	_, err := ts.GrantAccessToken(req)

	assert.EqualError(t, err, util.ErrInvalidGrant.Error())
}

func TestTokenService_GrantAccessToken_ShouldReturnTokens_WhenClientAppPingIsValid(t *testing.T) {
	ts := newTokenService()

	res, err := ts.GrantAccessToken(&TokenRequest{
		clientId:  test_data.CibaSession9.ClientId,
		authReqId: test_data.CibaSession9.AuthReqId,
	})

	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.NotNil(t, res.AccessToken)
	assert.NotNil(t, res.IdToken)
}
