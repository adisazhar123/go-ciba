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

	idToken := "eyJhbGciOiJSUzI1NiIsImtpZCI6IjNkNTg1Yzg4LWIyYWMtNGEwNy04MjRmLTY0OWNlYzI2MGFhNSIsInR5cCI6Imp3dCJ9.eyJhdF9oYXNoIjoiWkRnek1HSm1OVEExTnpBeE9UZzJPREl4TWpVM1pXSXdaREE0Wm1RMllnPT0iLCJhdWQiOiJlNjczOGMxMS0zNjYzLTQzMDEtYmExMi0yMzEyZWI5YTEzNDQiLCJhdXRoX3RpbWUiOjE1OTYzOTE1MzMsImV4cCI6MTU5NjM5NTEzMywiaWF0IjoxNTk2MzkxNTMzLCJpc3MiOiJjaWJhLXNlcnZlci5leGFtcGxlLmNvbSIsIm5vbmNlIjoiNzI3NTQ5OTktNWZlNS00OTZiLTkxODItOTA2ZjE1YTdiODZlIiwic3ViIjoiM2U0M2YwOGMtOTM0YS00ZjZjLTgwYjYtNmYyMDVhYmQ5NTRiIiwidXJuOm9wZW5pZDpwYXJhbXM6and0OmNsYWltOmF1dGhfcmVxX2lkIjoiZjE2OWFmZWUtOGEyYS00MWQ1LWJjZGEtNDUwZjk5OTQyZjBlIn0.n433lsZewflso3Ucg8yuO6EFTRxy5Gy26hJnUxayCbPAlyZ7pEBK7zvTjEXGF1VNGvtYgPoeLHWyUdwNclusBMlyU1xnLwXj0ctgHJ8U5N8KhIttR4W8tE2eJ_KkkanT6PqpN_dobApUzg0Vd48r-5a0PPTlfrT_41zkVrgMbxlIpPXKF2ysQD7C4T0ab7zYLCMpvMtZ5XXYxIaobzRxnf3iG4L-VJd6KDp7Me-N1VRPXZfuVqNmg_ID2W47R7VWLlubq8tWl908BacEiMIUmga8NTyUVc2KiwI0XxlaUeogVlfdJxmkLS43XBuDqU_NSlBfmD0UKT5rvQcmz4LNpcLGS3W-hH2tYM8jHi-WpqRXrEDbxvvxgI3YbahhjRx9hPpv9qpjhKiBN3kHtjHPJCspghUtH4UPYr6aFtXlk2Bd9MO9AiXjnbK4S6iSFl5oBZSCWe8NDqDF1q2piXnRmcqLf_XfiwgriYMDnnqbTX8tAPnkAr7RzXvF_ufmaX7viG0jvyf3GAcBIfcMyVfPjvMEwSRCL5aOcVNtFe6cxQl1IO3cH7o0FWTaQyRF2mpjW2If4CgQNn4Sax_rOLgzIvXo7zGo7gIuV8McR33_VSNtFUdquashV8Glv3-A_ijBWzI-Od7vid-R0ioXxXJVDhvozERz5ytgdFZATuFAo6U"

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
