package go_ciba

import (
	"testing"

	"github.com/adisazhar123/go-ciba/test_data"
	"github.com/adisazhar123/go-ciba/util"
	"github.com/stretchr/testify/assert"
)

func TestResourceServer_HandleResourceRequest_ShouldReturnErrInvalidToken(t *testing.T) {
	rs := &resourceServer{
		accessTokenRepo: test_data.NewAccessTokenVolatileRepository(),
		scopeUtil:       util.ScopeUtil{},
	}
	invalidToken := "4C6C584A-EB31-4E11-A3E7-EADBBB573E96"

	err := rs.HandleResourceRequest(&ResourceRequest{
		accessToken: invalidToken,
	}, "")

	assert.NotNil(t, err)
	assert.EqualError(t, err, util.ErrInvalidToken.Error())
}

func TestResourceServer_HandleResourceRequest_ShouldReturnErrInsufficientScope(t *testing.T) {
	rs := &resourceServer{
		accessTokenRepo: test_data.NewAccessTokenVolatileRepository(),
		scopeUtil:       util.ScopeUtil{},
	}
	token := "C59D9FBC-D8E4-4B8B-A95B-14F931EE1AB3"

	err := rs.HandleResourceRequest(&ResourceRequest{accessToken: token}, "payment:write")

	assert.NotNil(t, err)
	assert.EqualError(t, err, util.ErrInsufficientScope.Error())
}

func TestResourceServer_HandleResourceRequest_ShouldReturnErrInsufficientScope2(t *testing.T) {
	rs := &resourceServer{
		accessTokenRepo: test_data.NewAccessTokenVolatileRepository(),
		scopeUtil:       util.ScopeUtil{},
	}
	token := test_data.AccessTokenValid.Value

	err := rs.HandleResourceRequest(&ResourceRequest{accessToken: token}, "openid email profile chat:write payment:write")

	assert.NotNil(t, err)
	assert.EqualError(t, err, util.ErrInsufficientScope.Error())
}

func TestResourceServer_HandleResourceRequest_ShouldReturnErrInvalidTokenWhenExpired(t *testing.T) {
	rs := &resourceServer{
		accessTokenRepo: test_data.NewAccessTokenVolatileRepository(),
		scopeUtil:       util.ScopeUtil{},
	}
	token := test_data.AccessTokenExpired.Value

	err := rs.HandleResourceRequest(&ResourceRequest{accessToken: token}, "")

	assert.NotNil(t, err)
	assert.EqualError(t, err, util.ErrInvalidToken.Error())
}

func TestResourceServer_HandleResourceRequest_ShouldSucceedWhenTokenIsValid(t *testing.T) {
	rs := &resourceServer{
		accessTokenRepo: test_data.NewAccessTokenVolatileRepository(),
		scopeUtil:       util.ScopeUtil{},
	}
	token := test_data.AccessTokenValid.Value

	err := rs.HandleResourceRequest(&ResourceRequest{accessToken: token}, "")

	assert.Nil(t, err)
}

func TestResourceServer_HandleResourceRequest_ShouldSucceedWhenTokenIsValid_CustomScope(t *testing.T) {
	rs := &resourceServer{
		accessTokenRepo: test_data.NewAccessTokenVolatileRepository(),
		scopeUtil:       util.ScopeUtil{},
	}
	token := test_data.AccessTokenValid.Value

	err := rs.HandleResourceRequest(&ResourceRequest{accessToken: token}, "chat:write")

	assert.Nil(t, err)
}
