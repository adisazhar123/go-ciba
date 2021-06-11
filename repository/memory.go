package repository

import (
	"sync"

	"github.com/adisazhar123/go-ciba/domain"
)

type InMemoryDataStore struct {
	accessTokenRepo       *accessTokenRepository
	cibaSessionRepo       *cibaSessionRepository
	clientApplicationRepo *clientApplicationRepository
	keyRepo               *keyRepository
	userAccountRepo       *userAccountRepository
}

type accessTokenRepository struct {
	data map[string]*domain.AccessToken
	mtx  sync.RWMutex
}

func (a *accessTokenRepository) Create(accessToken *domain.AccessToken) error {
	a.mtx.Lock()
	defer a.mtx.Unlock()
	a.data[accessToken.Value] = accessToken
	return nil
}

func (a *accessTokenRepository) Find(accessToken string) (*domain.AccessToken, error) {
	a.mtx.RLock()
	defer a.mtx.RUnlock()
	return a.data[accessToken], nil
}

type cibaSessionRepository struct {
	data map[string]*domain.CibaSession
	mtx  sync.RWMutex
}

func (c *cibaSessionRepository) Create(cibaSession *domain.CibaSession) error {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	c.data[cibaSession.AuthReqId] = cibaSession
	return nil
}

func (c *cibaSessionRepository) FindById(id string) (*domain.CibaSession, error) {
	c.mtx.RLock()
	defer c.mtx.RUnlock()
	return c.data[id], nil
}

func (c *cibaSessionRepository) Update(cibaSession *domain.CibaSession) error {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	c.data[cibaSession.AuthReqId] = cibaSession
	return nil
}

type clientApplicationRepository struct {
	data map[string]*domain.ClientApplication
	mtx  sync.RWMutex
}

func (c *clientApplicationRepository) Register(clientApp *domain.ClientApplication) error {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	c.data[clientApp.Id] = clientApp
	return nil
}

func (c *clientApplicationRepository) FindById(id string) (*domain.ClientApplication, error) {
	c.mtx.RLock()
	defer c.mtx.RUnlock()
	return c.data[id], nil
}

type keyRepository struct {
	data map[string]*domain.Key
	mtx  sync.RWMutex
}

func (k *keyRepository) FindPrivateKeyByClientId(clientId string) (*domain.Key, error) {
	k.mtx.RLock()
	defer k.mtx.RUnlock()
	return k.data[clientId], nil
}

type userAccountRepository struct {
	data map[string]*domain.UserAccount
	mtx  sync.RWMutex
}

func (u *userAccountRepository) FindById(id string) (*domain.UserAccount, error) {
	u.mtx.RLock()
	defer u.mtx.RUnlock()
	return u.data[id], nil
}

func NewInMemoryDataStore() *InMemoryDataStore {
	return &InMemoryDataStore{
		accessTokenRepo: &accessTokenRepository{
			data: make(map[string]*domain.AccessToken),
			mtx:  sync.RWMutex{},
		},
		cibaSessionRepo: &cibaSessionRepository{
			data: make(map[string]*domain.CibaSession),
			mtx:  sync.RWMutex{},
		},
		clientApplicationRepo: &clientApplicationRepository{
			data: make(map[string]*domain.ClientApplication),
			mtx:  sync.RWMutex{},
		},
		keyRepo: &keyRepository{
			data: make(map[string]*domain.Key),
			mtx:  sync.RWMutex{},
		},
		userAccountRepo: &userAccountRepository{
			data: make(map[string]*domain.UserAccount),
			mtx:  sync.RWMutex{},
		},
	}
}

func (i *InMemoryDataStore) GetAccessTokenRepository() AccessTokenRepositoryInterface {
	return i.accessTokenRepo
}

func (i *InMemoryDataStore) GetCibaSessionRepository() CibaSessionRedisRepository {
	panic("implement me")
}

func (i *InMemoryDataStore) GetClientApplicationRepository() ClientApplicationRepositoryInterface {
	panic("implement me")
}

func (i *InMemoryDataStore) GetKeyRepository() KeyRepositoryInterface {
	panic("implement me")
}

func (i *InMemoryDataStore) GetUserAccountRepository() UserAccountRepositoryInterface {
	panic("implement me")
}
