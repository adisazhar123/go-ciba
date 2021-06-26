package util

import (
	"crypto/x509"
	"encoding/pem"
	"log"

	"gopkg.in/square/go-jose.v2"
	"gopkg.in/square/go-jose.v2/jwt"
)

type EncryptionInterface interface {
	Encode(payload interface{}, key, alg, keyId string) (string, error)
	Decode(jwt string, key string) (string, error)
}

type GoJoseEncryption struct {
}

func NewGoJoseEncryption() *GoJoseEncryption {
	return &GoJoseEncryption{}
}

func (gje *GoJoseEncryption) Decode(serialized string, key string) (string, error) {
	object, err := jose.ParseSigned(serialized)
	if err != nil {
		return "", err
	}
	output, err := object.Verify([]byte(key))
	if err != nil {
		return "", err
	}
	return string(output), nil
}

func (gje *GoJoseEncryption) Encode(payload interface{}, key, alg, keyId string) (string, error) {
	jwtKey := []byte(key)
	block, _ := pem.Decode(jwtKey)
	pKey, _ := x509.ParsePKCS1PrivateKey(block.Bytes)
	opt := &jose.SignerOptions{
		NonceSource: nil,
		EmbedJWK:    false,
	}
	opt.WithType("jwt")
	opt.WithBase64(true)
	opt.WithHeader("kid", keyId)

	sig, err := jose.NewSigner(jose.SigningKey{
		Algorithm: jose.SignatureAlgorithm(alg),
		Key:       pKey,
	}, opt)

	if err != nil {
		log.Printf("[go-ciba][encryption] an error occured create new signer: %s\n", err.Error())
		return "", err
	}

	raw, err := jwt.Signed(sig).Claims(payload).CompactSerialize()

	if err != nil {
		log.Printf("[go-ciba][encryption] an error occured serializing jwt claims: %s\n", err.Error())
		return "", err
	}
	return raw, nil
}
