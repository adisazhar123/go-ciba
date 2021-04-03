package repository

import (
	"context"
	"errors"
	"fmt"

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

func (c CibaSessionRedisRepository) FindById(id string) (*domain.CibaSession, error) {
	key := fmt.Sprintf("ciba_session:%s", id)
	val, err := c.client.Get(c.ctx, key).Result()
	if val == "" {
		return nil, errors.New("ciba session not found")
	}
	if err != nil {
		return nil, err
	}
	cibaSession := &domain.CibaSession{}
	if err := cibaSession.UnmarshalBinary([]byte(val)); err != nil {
		return nil, err
	}
	return cibaSession, nil
}

func (c CibaSessionRedisRepository) Update(cibaSession *domain.CibaSession) error {
	key := fmt.Sprintf("ciba_session:%s", cibaSession.AuthReqId)
	return c.client.Set(c.ctx, key, cibaSession, 0).Err()
}

type KeyRedisRepository struct {
	client *redis.Client
	ctx    context.Context
}

func NewKeyRedisRepository(addr string) *KeyRedisRepository {
	cli := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	return &KeyRedisRepository{
		client: cli,
		ctx:    context.Background(),
	}
}

func (k *KeyRedisRepository) FindPrivateKeyByClientId(clientId string) (*domain.Key, error) {
	key := fmt.Sprintf("oauth_key:%s", clientId)
	val, err := k.client.Get(k.ctx, key).Result()
	if val == "" {
		return nil, errors.New("oauth key not found")
	}
	if err != nil {
		return nil, err
	}
	oauthKey := &domain.Key{}
	if err := oauthKey.UnmarshalBinary([]byte(val)); err != nil {
		return nil, err
	}
	return oauthKey, nil
}
