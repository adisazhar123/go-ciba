package http_auth

import (
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"testing"

	"github.com/adisazhar123/go-ciba/domain"
	"github.com/adisazhar123/go-ciba/util"
	"github.com/stretchr/testify/assert"
)

func TestClientJwt_ValidateRequest_ShouldReturnFalseWhenClientAssertionIsMissing(t *testing.T) {
	formData := url.Values{}
	req, _ := http.NewRequest("POST", "/auth", strings.NewReader(formData.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(formData.Encode())))

	cJwt := clientJwt{goJose: util.NewGoJoseEncryption()}

	success := cJwt.ValidateRequest(req, &domain.ClientApplication{
		Id: "id_123",
	})

	assert.False(t, success)
}

func TestClientJwt_ValidateRequest_ShouldReturnFalseWhenClientAssertionTypeIsNotJwt(t *testing.T) {
	formData := url.Values{
		"client_assertion":      {"some.jwt.value"},
		"client_assertion_type": {"random_value"},
	}
	req, _ := http.NewRequest("POST", "/auth", strings.NewReader(formData.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(formData.Encode())))
	cJwt := clientJwt{goJose: util.NewGoJoseEncryption()}

	success := cJwt.ValidateRequest(req, &domain.ClientApplication{
		Id: "id_123",
	})

	assert.False(t, success)
}

func TestClientJwt_ValidateRequest_ShouldReturnFalseWhenClientSecretIsNotCorrect(t *testing.T) {
	formData := url.Values{
		// decoded value is "Lorem ipsum dolor sit amet"
		"client_assertion":      {"eyJhbGciOiJIUzI1NiJ9.TG9yZW0gaXBzdW0gZG9sb3Igc2l0IGFtZXQ.MRDTMykupahkRdvpsB8NSfgUrticeSSZ0kMiwyrLoZM"},
		"client_assertion_type": {jwtBearerAssertionType},
	}
	req, _ := http.NewRequest("POST", "/auth", strings.NewReader(formData.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(formData.Encode())))
	cJwt := clientJwt{goJose: util.NewGoJoseEncryption()}

	success := cJwt.ValidateRequest(req, &domain.ClientApplication{
		Id: "id_123",
	})

	assert.False(t, success)
}

func TestClientJwt_ValidateRequest_ShouldReturnFalseWhenRequiredClaimsAreMissingBecauseItIsNotAJsonString(t *testing.T) {
	formData := url.Values{
		// decoded value is "Lorem ipsum dolor sit amet"
		"client_assertion":      {"eyJhbGciOiJIUzI1NiJ9.TG9yZW0gaXBzdW0gZG9sb3Igc2l0IGFtZXQ.MRDTMykupahkRdvpsB8NSfgUrticeSSZ0kMiwyrLoZM"},
		"client_assertion_type": {jwtBearerAssertionType},
	}
	req, _ := http.NewRequest("POST", "/auth", strings.NewReader(formData.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(formData.Encode())))
	cJwt := clientJwt{goJose: util.NewGoJoseEncryption()}

	success := cJwt.ValidateRequest(req, &domain.ClientApplication{
		Id:     "id_123",
		Secret: "secret-key-123",
	})

	assert.False(t, success)
}

func TestClientJwt_ValidateRequest_ShouldReturnFalseWhenRequiredClaimsAreMissing(t *testing.T) {
	formData := url.Values{
		// decoded value is "{"iss":"id_123"}"
		"client_assertion":      {"eyJhbGciOiJIUzI1NiJ9.eyJpc3MiOiJpZF8xMjMifQ.ZsSQTJ3rpuKESkPMWTJyCgIvvyb-SpHCyNj_aP8O-Vs"},
		"client_assertion_type": {jwtBearerAssertionType},
	}
	req, _ := http.NewRequest("POST", "/auth", strings.NewReader(formData.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(formData.Encode())))
	cJwt := clientJwt{goJose: util.NewGoJoseEncryption()}

	success := cJwt.ValidateRequest(req, &domain.ClientApplication{
		Id:     "id_123",
		Secret: "secret-key-123",
	})

	assert.False(t, success)
}

func TestClientJwt_ValidateRequest_ShouldReturnTrueWhenGivenCorrectCredentials(t *testing.T) {
	formData := url.Values{
		// decoded value is "{"iss":"id_123","sub":"id_123","aud":"issuer-ciba.example.com/token","jti":"jti_123","exp":"2021-06-26T09:03:24.289326Z","iat":"2021-06-26T08:03:24.289326Z"}"
		"client_assertion":      {"eyJhbGciOiJIUzI1NiJ9.eyJpc3MiOiJpZF8xMjMiLCJzdWIiOiJpZF8xMjMiLCJhdWQiOiJpc3N1ZXItY2liYS5leGFtcGxlLmNvbS90b2tlbiIsImp0aSI6Imp0aV8xMjMiLCJleHAiOiIyMDIxLTA2LTI2VDA5OjAzOjI0LjI4OTMyNloiLCJpYXQiOiIyMDIxLTA2LTI2VDA4OjAzOjI0LjI4OTMyNloifQ.amWlTwUFNjh0JO6FNqKAGMH_pcgudQsbiUBrycf2Gd0"},
		"client_assertion_type": {jwtBearerAssertionType},
	}
	req, _ := http.NewRequest("POST", "/auth", strings.NewReader(formData.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(formData.Encode())))
	cJwt := clientJwt{goJose: util.NewGoJoseEncryption(), authServerTokenEndpoint: "issuer-ciba.example.com/token"}

	success := cJwt.ValidateRequest(req, &domain.ClientApplication{
		Id:     "id_123",
		Secret: "secret-key-123",
	})

	assert.True(t, success)
}
