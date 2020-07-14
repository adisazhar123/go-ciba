package util

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSliceStringContains_ContainsElement(t *testing.T) {
	slice := []string{"abc", "def", "ghi", "jkl"}
	find := "def"
	res := SliceStringContains(slice, find)

	expected := true

	assert.Equal(t, expected, res)
}

func TestSliceStringContains_DoesntContainsElement(t *testing.T) {
	slice := []string{"abc", "def", "ghi", "jkl"}
	find := "klmn"
	res := SliceStringContains(slice, find)

	expected := false

	assert.Equal(t, expected, res)
}
