package util

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGenerateRandomString(t *testing.T) {
	string := GenerateRandomString()
	assert.NotEmpty(t, string)
}

func TestGenerateUuid(t *testing.T) {
	uuid := GenerateUuid()
	assert.NotEmpty(t, uuid)
}
