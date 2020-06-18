package domain

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestNewCibaSession(t *testing.T) {
	cs := NewCibaSession(120, 5)

	assert.Equal(t, 120, cs.expiresIn,)
	assert.Equal(t, 5, cs.interval)
	assert.NotEmpty(t, cs.authReqId)

	assert.NotEqual(t, 50, cs.expiresIn)
	assert.NotEqual(t, 13, cs.interval)
}