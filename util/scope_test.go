package util

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestScopeUtil_ScopeExist_ValidScope(t *testing.T) {
	util := ScopeUtil{}
	requestedScope := "openid profile"
	registeredScope := "openid profile address"

	requestedScope2 := registeredScope
	registeredScope2 := "openid profile address"

	res := util.ScopeExist(registeredScope, requestedScope)
	res2 := util.ScopeExist(registeredScope2, requestedScope2)

	assert.Equal(t, true, res)
	assert.Equal(t, true, res2)
}

func TestScopeUtil_ScopeExist_InvalidScope(t *testing.T) {
	util := ScopeUtil{}
	requestedScope := "wrong-scope openid"
	registeredScope := "openid profile address email"

	requestedScope2 := "wrong-scope openid profile address email"
	registeredScope2 := registeredScope

	res := util.ScopeExist(registeredScope, requestedScope)
	res2 := util.ScopeExist(registeredScope2, requestedScope2)

	assert.Equal(t, false, res)
	assert.Equal(t, false, res2)
}
