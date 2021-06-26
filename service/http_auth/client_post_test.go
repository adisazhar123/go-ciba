package http_auth

import (
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"testing"

	"github.com/adisazhar123/go-ciba/domain"
	"github.com/stretchr/testify/assert"
)

func TestClientPost_ValidateRequest_ShouldReturnTrueWhenGivenCorrectCredentials(t *testing.T) {
	formData := url.Values{
		"client_id": {"id_123"},
		"client_secret": {"secret_123"},
	}
	req, _ := http.NewRequest("POST", "/auth", strings.NewReader(formData.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(formData.Encode())))
	cp := &clientPost{}

	success := cp.ValidateRequest(req, &domain.ClientApplication{
		Id:     "id_123",
		Secret: "secret_123",
	})

	assert.True(t, success)
}

func TestClientPost_ValidateRequest_ShouldReturnFalseWhenGivenIncorrectCredentials(t *testing.T) {
	formData := url.Values{
		"client_id": {"id_123_234234"},
		"client_secret": {"secret_123"},
	}
	req, _ := http.NewRequest("POST", "/auth", strings.NewReader(formData.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(formData.Encode())))
	cp := &clientPost{}

	success := cp.ValidateRequest(req, &domain.ClientApplication{
		Id:     "id_123",
		Secret: "secret_123",
	})

	assert.False(t, success)
}

func TestClientPost_ValidateRequest_ShouldReturnFalseWhenRequestIsMissingClientId(t *testing.T) {
	formData := url.Values{
		"client_secret": {"secret_123"},
	}
	req, _ := http.NewRequest("POST", "/auth", strings.NewReader(formData.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(formData.Encode())))
	cp := &clientPost{}

	success := cp.ValidateRequest(req, &domain.ClientApplication{
		Id:     "id_123",
		Secret: "secret_123",
	})

	assert.False(t, success)
}

func TestClientPost_ValidateRequest_ShouldReturnFalseWhenRequestIsMissingClientSecret(t *testing.T) {
	formData := url.Values{
		"client_id": {"id_123"},
	}
	req, _ := http.NewRequest("POST", "/auth", strings.NewReader(formData.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(formData.Encode())))
	cp := &clientPost{}

	success := cp.ValidateRequest(req, &domain.ClientApplication{
		Id:     "id_123",
		Secret: "secret_123",
	})

	assert.False(t, success)
}

func TestClientPost_ValidateRequest_ShouldReturnFalseWhenRequestIsMissingClientIdAndClientSecret(t *testing.T) {
	formData := url.Values{}
	req, _ := http.NewRequest("POST", "/auth", strings.NewReader(formData.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(formData.Encode())))
	cp := &clientPost{}

	success := cp.ValidateRequest(req, &domain.ClientApplication{
		Id:     "id_123",
		Secret: "secret_123",
	})

	assert.False(t, success)
}