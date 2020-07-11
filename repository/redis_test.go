package repository

import (
	"github.com/adisazhar123/ciba-server/domain"
	"github.com/adisazhar123/ciba-server/util"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestClientApplicationRedisRepository_Register(t *testing.T) {
	repo := NewClientApplicationRedisRepository()
	randomId := util.GenerateUuid()
	newClientApp := domain.NewClientApplication("test-app-" + randomId, "openid profile email", "ping", "https://adisazhar.com/notification", "RS256", false)

	err := repo.Register(newClientApp)
	assert.Empty(t, err)
}
