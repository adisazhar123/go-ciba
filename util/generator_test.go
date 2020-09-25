package util

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGenerateRandomString(t *testing.T) {
	randomString := GenerateRandomString()
	assert.NotEmpty(t, randomString)
}

func TestGenerateUuid(t *testing.T) {
	uuid := GenerateUuid()
	assert.NotEmpty(t, uuid)
}
