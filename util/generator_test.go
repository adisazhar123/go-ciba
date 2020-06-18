package util

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGenerateUuid(t *testing.T) {
	uuid := GenerateUuid()
	assert.NotEmpty(t, uuid)
}