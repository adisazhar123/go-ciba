package domain

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewCibaSession(t *testing.T) {
	cs := NewCibaSession(120, 5)

	assert.Equal(t, 120, cs.expiresIn,)
	assert.Equal(t, 5, cs.interval)
	assert.NotEmpty(t, cs.authReqId)

	assert.Greater(t, cs.expiresIn, 0)
	assert.Greater(t, cs.interval, 0)

	assert.NotEqual(t, 50, cs.expiresIn)
	assert.NotEqual(t, 13, cs.interval)
}