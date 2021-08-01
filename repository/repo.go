package repository

import "github.com/adisazhar123/go-ciba/domain"

type common interface {
	HaveTransactionSupport() bool
}

type AccessTokenRepositoryInterface interface {
	Create(accessToken *domain.AccessToken) error
	Find(accessToken string) (*domain.AccessToken, error)
}

type CibaSessionRepositoryInterface interface {
	Create(cibaSession *domain.CibaSession) error
	FindById(id string) (*domain.CibaSession, error)
	Update(cibaSession *domain.CibaSession) error
}

type ClientApplicationRepositoryInterface interface {
	Register(clientApp *domain.ClientApplication) error
	FindById(id string) (*domain.ClientApplication, error)
}

type KeyRepositoryInterface interface {
	FindPrivateKeyByClientId(clientId string) (*domain.Key, error)
}

type UserAccountRepositoryInterface interface {
	FindById(id string) (*domain.UserAccount, error)
}

type UserClaimRepositoryInterface interface {
	GetUserClaims(userId, scopes string) (map[string]interface{}, error)
}

type DataStoreInterface interface {
	GetAccessTokenRepository() AccessTokenRepositoryInterface
	GetCibaSessionRepository() CibaSessionRepositoryInterface
	GetClientApplicationRepository() ClientApplicationRepositoryInterface
	GetKeyRepository() KeyRepositoryInterface
	GetUserAccountRepository() UserAccountRepositoryInterface
	GetUserClaimRepository() UserClaimRepositoryInterface
}
