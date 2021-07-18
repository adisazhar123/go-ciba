package repository

import "github.com/adisazhar123/go-ciba/domain"

type common interface {
	HaveTransactionSupport() bool
}

// to store access token.
type AccessTokenRepositoryInterface interface {
	// common
	Create(accessToken *domain.AccessToken) error
	Find(accessToken string) (*domain.AccessToken, error)
}

// to store ciba session.
type CibaSessionRepositoryInterface interface {
	// common
	Create(cibaSession *domain.CibaSession) error
	FindById(id string) (*domain.CibaSession, error)
	Update(cibaSession *domain.CibaSession) error
}

// to store client application.
type ClientApplicationRepositoryInterface interface {
	// common
	Register(clientApp *domain.ClientApplication) error
	FindById(id string) (*domain.ClientApplication, error)
}

// to store public & private key.
type KeyRepositoryInterface interface {
	// common
	FindPrivateKeyByClientId(clientId string) (*domain.Key, error)
}

type UserAccountRepositoryInterface interface {
	// common
	FindById(id string) (*domain.UserAccount, error)
}

type UserClaimRepositoryInterface interface {
	GetUserClaims(userId, scopes string) map[string]interface{}
}

type DataStoreInterface interface {
	GetAccessTokenRepository() AccessTokenRepositoryInterface
	GetCibaSessionRepository() CibaSessionRepositoryInterface
	GetClientApplicationRepository() ClientApplicationRepositoryInterface
	GetKeyRepository() KeyRepositoryInterface
	GetUserAccountRepository() UserAccountRepositoryInterface
	GetUserClaimRepository() UserClaimRepositoryInterface
}
