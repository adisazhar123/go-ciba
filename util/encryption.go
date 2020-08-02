package util

import (
	"crypto/x509"
	"encoding/pem"
	"gopkg.in/square/go-jose.v2"
	"gopkg.in/square/go-jose.v2/jwt"
)

type EncryptionInterface interface {
	Encode(payload interface{}, key, alg, keyId string) string
	Decode(jwt interface{}, key string, allowedAlgorithms string) interface{}
}

type GoJoseEncryption struct {
}

func (gje *GoJoseEncryption) Decode(jwt interface{}, key string, allowedAlgorithms string) interface{} {
	panic("implement me")
}

func (gje *GoJoseEncryption) Encode(payload interface{}, key, alg, keyId string) string {
	jwtKey := []byte(key)
	block, _ := pem.Decode(jwtKey)
	pKey, _ := x509.ParsePKCS1PrivateKey(block.Bytes)

	sig, err := jose.NewSigner(jose.SigningKey{
		Algorithm: jose.SignatureAlgorithm(alg),
		Key:       pKey,
	}, &jose.SignerOptions{
		NonceSource: nil,
		EmbedJWK:    false,
		ExtraHeaders: map[jose.HeaderKey]interface{}{
			"kid": keyId,
			"typ": "jwt",
		},
	})

	if err != nil {
		panic(err)
	}

	raw, err := jwt.Signed(sig).Claims(payload).CompactSerialize()

	if err != nil {
		panic(err)
	}
	return raw
}
