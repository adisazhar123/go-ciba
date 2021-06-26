package util

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGoJoseEncryption_Encode_ShouldFailBecauseOfIncorrectKey(t *testing.T) {
	enc := NewGoJoseEncryption()
	payload, err := enc.Decode("eyJhbGciOiJIUzI1NiJ9.TG9yZW0gaXBzdW0gZG9sb3Igc2l0IGFtZXQ.MRDTMykupahkRdvpsB8NSfgUrticeSSZ0kMiwyrLoZM", "111111")

	assert.Error(t, err)
	assert.EqualError(t, err, "square/go-jose: error in cryptographic primitive")
	assert.Empty(t, payload)
}

func TestGoJoseEncryption_Encode_ShouldSucceed(t *testing.T) {
	enc := NewGoJoseEncryption()
	payload, err := enc.Decode("eyJhbGciOiJIUzI1NiJ9.TG9yZW0gaXBzdW0gZG9sb3Igc2l0IGFtZXQ.MRDTMykupahkRdvpsB8NSfgUrticeSSZ0kMiwyrLoZM", "secret-key-123")

	assert.NoError(t, err)
	assert.Equal(t, "Lorem ipsum dolor sit amet", payload)
}