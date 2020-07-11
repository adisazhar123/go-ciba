package repository

import (
	"context"

	"github.com/adisazhar123/ciba-server/domain"
	"github.com/go-redis/redis/v8"
)

func NewRedisClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
}

func NewClientApplicationRedisRepository() *ClientApplicationRedisRepository {
	return &ClientApplicationRedisRepository{
		client: NewRedisClient(),
		ctx:    context.Background(),
	}
}

type ClientApplicationRedisRepository struct {
	client *redis.Client
	ctx context.Context
}

func (ca *ClientApplicationRedisRepository) Register(clientApp *domain.ClientApplication) error {
	return ca.client.Set(ca.ctx, "client_application:" + clientApp.GetId(), clientApp, 0).Err()
}

type UserAccountRedisRepository struct {
	client *redis.Client
	ctx context.Context
}

func (ua *UserAccountRedisRepository) FindById(id string) (*domain.UserAccount, error) {
	val, err := ua.client.Get(ua.ctx, "user_account:" + id).Result()
	if err != nil {
		return nil, err
	}

	user := &domain.UserAccount{}
	if err := user.UnmarshalBinary([]byte(val)); err != nil {
		return nil, err
	}

	return user, nil
}