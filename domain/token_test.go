package domain

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/adisazhar123/go-ciba/util"
	"github.com/stretchr/testify/assert"
)

func newIdTokenManager() TokenManager {
	return TokenManager{e: newEncryptionUtil()}
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

	idToken := "eyJhbGciOiJSUzI1NiIsImtpZCI6IjNkNTg1Yzg4LWIyYWMtNGEwNy04MjRmLTY0OWNlYzI2MGFhNSIsInR5cCI6Imp3dCJ9.eyJhdF9oYXNoIjoiWkRnek1HSm1OVEExTnpBeE9UZzJPREl4TWpVM1pXSXdaREE0Wm1RMllnIiwiYXVkIjoiZTY3MzhjMTEtMzY2My00MzAxLWJhMTItMjMxMmViOWExMzQ0IiwiYXV0aF9yZXFfaWQiOiJmMTY5YWZlZS04YTJhLTQxZDUtYmNkYS00NTBmOTk5NDJmMGUiLCJhdXRoX3RpbWUiOjE1OTYzOTE1MzMsImV4cCI6MTU5NjM5NTEzMywiaWF0IjoxNTk2MzkxNTMzLCJpc3MiOiJjaWJhLXNlcnZlci5leGFtcGxlLmNvbSIsIm5vbmNlIjoiNzI3NTQ5OTktNWZlNS00OTZiLTkxODItOTA2ZjE1YTdiODZlIiwic3ViIjoiM2U0M2YwOGMtOTM0YS00ZjZjLTgwYjYtNmYyMDVhYmQ5NTRiIn0.sSU0SQvGEBEzIM-Iqp5m7mz69izDmxr-K7jd7UYGPwpBcZlzvuElAkjLE6dE5-chSjGLNtkd8F4uE864fYTV6L8CM1BzhrPAj-XOTCH2nSIJfytrIVCoOm3r_jF-Oyo49XpzDgZOoagoZXd6vwMh8HqKUMw-mNvGWsCiKmVLdy2yTFEt3F1fUGweXSdzbrHRMK6jeWNqBaM8nR3vt-Z2ddy0_4zV2yCdEoKKswLGJfTEnmNRITx5JO_vkAn_vtVX2DBCuVz_C-YBBtZhP0AmsfSLeqMm-oXEvvoh4dzkHddWH27VnylgGO1pJh5On68xSoP9xiGzLQpxRjSHw2I5zUFRsMqtBmorYaCq25CzoAuLjnORECUv8uwBjIJ_AR4FPkWVZLWcwXnIUaEZv1Z2fldBqLpVYKZoIIQ996DNbwHmw6gvfxPZw623CgfRIJZ_FHKW4HCyjt0dxuzShM3NN509w7TBcBuS4CGzi0iwRenSp9A1mo_iIgQdmpOX8IobHYNL1vykSbZeUSGIKEwwzupHpNs2JoLY-7MO-cll0IoHdf9m2KNref4-qBmEKwJDtFaKnx4z0PjPd3rzTHP0YFBusx2NMNltbBurdf8qZ0vzGOx5y2PiH6Y1iz2z_qglRkNJjRlysaI5jr9Z1XwA5COrCcaKfFKVpvPG9_Gb7aY"

	claims := make(map[string]interface{})
	claims["auth_req_id"] = "f169afee-8a2a-41d5-bcda-450f99942f0e"
	claims["aud"] = "e6738c11-3663-4301-ba12-2312eb9a1344"
	claims["auth_time"] = now
	claims["iat"] = now
	claims["exp"] = now + 3600
	claims["iss"] = "ciba-server.example.com"
	claims["sub"] = "3e43f08c-934a-4f6c-80b6-6f205abd954b"
	claims["nonce"] = "72754999-5fe5-496b-9182-906f15a7b86e"

	et := mgr.CreateIdToken(claims, string(privateKey), alg, kId, token)

	assert.NotNil(t, et)
	assert.Equal(t, idToken, et.Value)
}
