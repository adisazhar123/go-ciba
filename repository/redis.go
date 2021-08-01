package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/adisazhar123/go-ciba/domain"
	"github.com/go-redis/redis/v8"
)

func NewClientApplicationRedisRepository(client *redis.Client) *clientApplicationRedisRepository {
	return &clientApplicationRedisRepository{
		client: client,
		ctx:    context.Background(),
	}
}

type clientApplicationRedisRepository struct {
	client *redis.Client
	ctx    context.Context
}

func (ca *clientApplicationRedisRepository) FindById(id string) (*domain.ClientApplication, error) {
	key := fmt.Sprintf("client_application:%s", id)
	val, err := ca.client.Get(ca.ctx, key).Result()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	clientApp := &domain.ClientApplication{}
	if err := clientApp.UnmarshalBinary([]byte(val)); err != nil {
		return nil, err
	}
	return clientApp, nil
}

func (ca *clientApplicationRedisRepository) Register(clientApp *domain.ClientApplication) error {
	key := fmt.Sprintf("client_application:%s", clientApp.GetId())
	return ca.client.Set(ca.ctx, key, clientApp, 0).Err()
}

type userAccountRedisRepository struct {
	client redis.Cmdable
	ctx    context.Context
}

func NewUserAccountRedisRepository(client *redis.Client) *userAccountRedisRepository {
	return &userAccountRedisRepository{
		client: client,
		ctx:    context.Background(),
	}
}

func (ua *userAccountRedisRepository) FindById(id string) (*domain.UserAccount, error) {
	key := fmt.Sprintf("user_account:%s", id)
	val, err := ua.client.Get(ua.ctx, key).Result()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	user := &domain.UserAccount{}
	if err := user.UnmarshalBinary([]byte(val)); err != nil {
		return nil, err
	}

	return user, nil
}

func NewCibaSessionRedisRepository(client *redis.Client) *CibaSessionRedisRepository {
	return &CibaSessionRedisRepository{
		client: client,
		ctx:    context.Background(),
	}
}

type CibaSessionRedisRepository struct {
	client *redis.Client
	ctx    context.Context
}

func (c *CibaSessionRedisRepository) Create(cibaSession *domain.CibaSession) error {
	key := fmt.Sprintf("ciba_session:%s", cibaSession.AuthReqId)
	return c.client.Set(c.ctx, key, cibaSession, 0).Err()
}

func (c *CibaSessionRedisRepository) FindById(id string) (*domain.CibaSession, error) {
	key := fmt.Sprintf("ciba_session:%s", id)
	val, err := c.client.Get(c.ctx, key).Result()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	cibaSession := &domain.CibaSession{}
	if err := cibaSession.UnmarshalBinary([]byte(val)); err != nil {
		return nil, err
	}
	return cibaSession, nil
}

func (c *CibaSessionRedisRepository) Update(cibaSession *domain.CibaSession) error {
	return c.Create(cibaSession)
}

type keyRedisRepository struct {
	client *redis.Client
	ctx    context.Context
}

func NewKeyRedisRepository(client *redis.Client) *keyRedisRepository {
	return &keyRedisRepository{
		client: client,
		ctx:    context.Background(),
	}
}

func (k *keyRedisRepository) FindPrivateKeyByClientId(clientId string) (*domain.Key, error) {
	key := fmt.Sprintf("oauth_key:%s", clientId)
	val, err := k.client.Get(k.ctx, key).Result()
	if err == redis.Nil {
		return nil, nil
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

type accessTokenRedisRepository struct {
	client *redis.Client
	ctx    context.Context
}

func NewAccessTokenRedisRepository(client *redis.Client) *accessTokenRedisRepository {
	return &accessTokenRedisRepository{
		client: client,
		ctx:    context.Background(),
	}
}

func (a *accessTokenRedisRepository) Create(accessToken *domain.AccessToken) error {
	key := fmt.Sprintf("access_token:%s", accessToken.Value)
	return a.client.Set(a.ctx, key, accessToken, 0).Err()
}

func (a *accessTokenRedisRepository) Find(accessToken string) (*domain.AccessToken, error) {
	key := fmt.Sprintf("access_token:%s", accessToken)
	val, err := a.client.Get(a.ctx, key).Result()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	at := &domain.AccessToken{}
	if err := at.UnmarshalBinary([]byte(val)); err != nil {
		return nil, err
	}
	return at, nil
}

type userClaimRedisRepository struct {
	client *redis.Client
	ctx    context.Context
}

func NewUserClaimRedisRepository(client *redis.Client) *userClaimRedisRepository {
	return &userClaimRedisRepository{
		client: client,
		ctx:    context.Background(),
	}
}

func (u *userClaimRedisRepository) GetUserClaims(userId, scopes string) (map[string]interface{}, error) {
	userKey := fmt.Sprintf("user_account:%s", userId)
	val, err := u.client.Get(u.ctx, userKey).Result()
	if err == redis.Nil {
		return map[string]interface{}{}, nil
	} else if err != nil {
		return nil, err
	}

	var userAccount map[string]interface{}
	if err := json.Unmarshal([]byte(val), &userAccount); err != nil {
		return nil, err
	}
	var claims []string
	scopesArr := strings.Split(scopes, " ")

	for _, scope := range scopesArr {
		claimsInScope, err := u.client.LRange(u.ctx, "scope:"+scope, 0, -1).Result()
		if err == redis.Nil {
			continue
		} else if err != nil {
			return nil, err
		}
		claims = append(claims, claimsInScope...)
	}
	claimsValues := make(map[string]interface{})
	for _, claim := range claims {
		val, ok := userAccount[claim]
		if ok {
			claimsValues[claim] = val
		}
	}
	return claimsValues, nil
}

type RedisDataStore struct {
	accessTokenRepo       *accessTokenRedisRepository
	cibaSessionRepo       *CibaSessionRedisRepository
	clientApplicationRepo *clientApplicationRedisRepository
	keyRepositoryRepo     *keyRedisRepository
	userAccountRepo       *userAccountRedisRepository
	userClaimRepo         *userClaimRedisRepository
}

func NewRedisDataStore(client *redis.Client) *RedisDataStore {
	return &RedisDataStore{
		accessTokenRepo:       NewAccessTokenRedisRepository(client),
		cibaSessionRepo:       NewCibaSessionRedisRepository(client),
		clientApplicationRepo: NewClientApplicationRedisRepository(client),
		keyRepositoryRepo:     NewKeyRedisRepository(client),
		userAccountRepo:       NewUserAccountRedisRepository(client),
		userClaimRepo:         NewUserClaimRedisRepository(client),
	}
}

func (r *RedisDataStore) GetAccessTokenRepository() AccessTokenRepositoryInterface {
	return r.accessTokenRepo
}

func (r *RedisDataStore) GetCibaSessionRepository() CibaSessionRepositoryInterface {
	return r.cibaSessionRepo
}

func (r *RedisDataStore) GetClientApplicationRepository() ClientApplicationRepositoryInterface {
	return r.clientApplicationRepo
}

func (r *RedisDataStore) GetKeyRepository() KeyRepositoryInterface {
	return r.keyRepositoryRepo
}

func (r *RedisDataStore) GetUserAccountRepository() UserAccountRepositoryInterface {
	return r.userAccountRepo
}

func (r *RedisDataStore) GetUserClaimRepository() UserClaimRepositoryInterface {
	return r.userClaimRepo
}
