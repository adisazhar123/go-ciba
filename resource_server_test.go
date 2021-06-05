package go_ciba

import (
	"testing"

	"github.com/adisazhar123/go-ciba/test_data"
	"github.com/adisazhar123/go-ciba/util"
	"github.com/stretchr/testify/assert"
)

func TestResourceServer_HandleResourceRequest_ShouldReturnErrInvalidToken(t *testing.T) {
	rs := &ResourceServer{
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
