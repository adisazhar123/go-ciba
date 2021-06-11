package repository

import "github.com/adisazhar123/go-ciba/domain"

// to store access token.
type AccessTokenRepositoryInterface interface {
	Create(accessToken *domain.AccessToken) error
	Find(accessToken string) (*domain.AccessToken, error)
}

// to store ciba session.
type CibaSessionRepositoryInterface interface {
	Create(cibaSession *domain.CibaSession) error
	FindById(id string) (*domain.CibaSession, error)
	Update(cibaSession *domain.CibaSession) error
}

// to store client application.
type ClientApplicationRepositoryInterface interface {
	Register(clientApp *domain.ClientApplication) error
	FindById(id string) (*domain.ClientApplication, error)
}

// to store public & private key.
type KeyRepositoryInterface interface {
	FindPrivateKeyByClientId(clientId string) (*domain.Key, error)
}

type UserAccountRepositoryInterface interface {
	FindById(id string) (*domain.UserAccount, error)
}

type DataStoreInterface interface {
	GetAccessTokenRepository() AccessTokenRepositoryInterface
	GetCibaSessionRepository() CibaSessionRedisRepository
	GetClientApplicationRepository() ClientApplicationRepositoryInterface
	GetKeyRepository() KeyRepositoryInterface
	GetUserAccountRepository() UserAccountRepositoryInterface
}
