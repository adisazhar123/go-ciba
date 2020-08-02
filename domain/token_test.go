package domain

import (
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"fmt"
	"github.com/adisazhar123/go-ciba/util"
	"github.com/stretchr/testify/assert"
	"hash"
	"io/ioutil"
	"os"
	"testing"
)

func newIdTokenManager() IdTokenManager {
	return IdTokenManager{e: newEncryptionUtil()}
}

func newEncryptionUtil() util.EncryptionInterface {
	return &util.GoJoseEncryption{}
}

func TestIdTokenManager_CreateIdToken(t *testing.T) {
	mgr := newIdTokenManager()
	file, _ := os.Open("../test_data/key.pem")
	defer file.Close()

	privateKey, _ := ioutil.ReadAll(file)
	alg := "RS256"
	kId := "3d585c88-b2ac-4a07-824f-649cec260aa5"
	token := "2b7d56d2-7f1e-407f-896d-71bebfbdc0d4"
	now := 1596391533

	defaultClaims := DefaultCibaIdTokenClaims{}

	defaultClaims.AuthReqId = "f169afee-8a2a-41d5-bcda-450f99942f0e"
	defaultClaims.Aud = "e6738c11-3663-4301-ba12-2312eb9a1344"
	defaultClaims.AuthTime = now
	defaultClaims.Iat = now
	defaultClaims.Exp = now + 3600
	defaultClaims.Iss = "ciba-server.example.com"
	defaultClaims.Sub = "3e43f08c-934a-4f6c-80b6-6f205abd954b"
	defaultClaims.AtHash = createTokenHash(token, alg)
	defaultClaims.Nonce = "72754999-5fe5-496b-9182-906f15a7b86e"

	idToken := "eyJhbGciOiJSUzI1NiIsImtpZCI6IjNkNTg1Yzg4LWIyYWMtNGEwNy04MjRmLTY0OWNlYzI2MGFhNSIsInR5cCI6Imp3dCJ9.eyJhdF9oYXNoIjoiWkRnek1HSm1OVEExTnpBeE9UZzJPREl4TWpVM1pXSXdaREE0Wm1RMllnPT0iLCJhdWQiOiJlNjczOGMxMS0zNjYzLTQzMDEtYmExMi0yMzEyZWI5YTEzNDQiLCJhdXRoX3RpbWUiOjE1OTYzOTE1MzMsImV4cCI6MTU5NjM5NTEzMywiaWF0IjoxNTk2MzkxNTMzLCJpc3MiOiJjaWJhLXNlcnZlci5leGFtcGxlLmNvbSIsIm5vbmNlIjoiNzI3NTQ5OTktNWZlNS00OTZiLTkxODItOTA2ZjE1YTdiODZlIiwic3ViIjoiM2U0M2YwOGMtOTM0YS00ZjZjLTgwYjYtNmYyMDVhYmQ5NTRiIiwidXJuOm9wZW5pZDpwYXJhbXM6and0OmNsYWltOmF1dGhfcmVxX2lkIjoiZjE2OWFmZWUtOGEyYS00MWQ1LWJjZGEtNDUwZjk5OTQyZjBlIn0.X04-EAgs-MrMY2zRHCmgGP3tOwOYQ-25u4dhp2v0_5LLMxGGHirpUwzh-tGCQXJM1M4YbYxo_YlDmbvkvLK7yX-N3Q6L4UfR0Dds_OUlBOCqY50H5RtbkzNZ7ZG31oRs_01crl4FTuZedUUX8VKWDMTlDg9BON0SyfdYYhgn2L38_uyupBDKmqVMql8IxLlRlJ2VQgjmqLsYsc7ag6XipcbeunLpBXlYgCL4rIrI-MNnLH3egozCjbQ9zDKXfnzgjPK2q9FaEWx0AL2JdQsSIJCb7KB6-74yefyzZhLGVW0mSSv70JM56gQY4q7p-m13XxLr2hwcpJiW5O86Y4RjMoDILbUXw8b9fwv5QTLpb2B8IxtPaoQUXGneyvRFTRr5AeL3KP5qvJVsbluv5r1CH3P6FWa63s9fwD5k_WNXx1GNRuXl_vhRZQeueH14F4LFA8GsHfxiTIMdj8QXDLj3GFhBa7LQ5cWtYQCOt3DwSezd7_n54m6ug6TRHP7mQOy09nbL9fGOkheVHZoZYqzg4B4KX5B4FCbOxJ8EQaUfM6oFUKEchm9n3lEts0KSH793-8uQgvs6tYJq_1B7LuxmHCxcD6mfyjrz8E1Le24YO5biTF7GNwUfKWE2QNY70wmRWOmNPlVe7zd-Zw9fr2wNQP_9eTZjkXh6AkFI8WoUc-8"

	et := mgr.CreateIdToken(defaultClaims, string(privateKey), alg, kId)

	assert.NotNil(t, et)
	assert.Equal(t, idToken, et.value)
}

func createTokenHash(token, alg string) string {
	alg = alg[2:]
	hashAlg := fmt.Sprintf("SHA%s", alg)

	var h hash.Hash

	if hashAlg == "SHA256" {
		h = sha256.New()
	} else if hashAlg == "SHA512" {
		h = sha512.New()
	} else {
		panic("hash algorithm not supported")
	}

	h.Write([]byte(token))
	hashed := h.Sum(nil)
	hashStr := fmt.Sprintf("%x", hashed)
	tokenHash := hashStr[:(len(hashStr)/2)-1]

	return base64.URLEncoding.EncodeToString([]byte(tokenHash))
}
