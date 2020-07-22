package repository

import (
	"context"
	"github.com/adisazhar123/go-ciba/domain"
	"github.com/go-redis/redis/v8"
)

func NewClientApplicationRedisRepository(addr string) *ClientApplicationRedisRepository {
	cli := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	return &ClientApplicationRedisRepository{
		client: cli,
		ctx:    context.Background(),
	}
}

type ClientApplicationRedisRepository struct {
	client *redis.Client
	ctx    context.Context
}

func (ca *ClientApplicationRedisRepository) Register(clientApp *domain.ClientApplication) error {
	return ca.client.Set(ca.ctx, "client_application:"+clientApp.GetId(), clientApp, 0).Err()
}

type UserAccountRedisRepository struct {
	client redis.Cmdable
	ctx    context.Context
}

func NewUserAccountRedisRepository(addr string) *UserAccountRedisRepository {
	cli := redis.NewClient(&redis.Options{
		Addr: addr,
	})
	return &UserAccountRedisRepository{
		client: cli,
		ctx:    context.Background(),
	}
}

func (ua *UserAccountRedisRepository) FindById(id string) (*domain.UserAccount, error) {
	val, err := ua.client.Get(ua.ctx, "user_account:"+id).Result()
	if val == "" || err != nil {
		return nil, err
	}

	user := &domain.UserAccount{}
	if err := user.UnmarshalBinary([]byte(val)); err != nil {
		return nil, err
	}

	return user, nil
}

func NewCibaSessionRedisRepository(addr string) *CibaSessionRedisRepository {
	cli := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	return &CibaSessionRedisRepository{
		client: cli,
		ctx:    context.Background(),
	}
}

type CibaSessionRedisRepository struct {
	client *redis.Client
	ctx    context.Context
}

func (c CibaSessionRedisRepository) Create(cibaSession *domain.CibaSession) error {
	return c.client.Set(c.ctx, "ciba_session:"+cibaSession.AuthReqId, cibaSession, 0).Err()
}
