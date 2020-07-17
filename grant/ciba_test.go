package grant

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCibaGrant_GetIdentifier(t *testing.T) {
	ciba := NewCibaGrant()
	id := "urn:openid:params:grant-type:ciba"

	assert.Equal(t, id, ciba.GetIdentifier())
}

func TestNewCibaGrant(t *testing.T) {
	ciba := NewCibaGrant()

	assert.NotNil(t, ciba)
}
