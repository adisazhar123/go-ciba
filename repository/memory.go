package repository

import (
	"io/ioutil"
	"os"
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
	privateKeyFile, _ := os.Open("../data/key.pem")
	publicKeyFile, _ := os.Open("../data/public.pem")

	privateKey, _ := ioutil.ReadAll(privateKeyFile)
	publicKey, _ := ioutil.ReadAll(publicKeyFile)

	clientId := "4E6E4BE4-089F-40C1-A4EE-BE40CE119AAE"

	key := domain.Key{
		ID:       "1",
		ClientId: "unknown",
		Alg:      "RS256",
		Public:   string(publicKey),
		Private:  string(privateKey),
	}

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
			data: map[string]*domain.Key{
				clientId: &key,
			},
			mtx: sync.RWMutex{},
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

func (i *InMemoryDataStore) GetCibaSessionRepository() CibaSessionRepositoryInterface {
	return i.cibaSessionRepo
}

func (i *InMemoryDataStore) GetClientApplicationRepository() ClientApplicationRepositoryInterface {
	return i.clientApplicationRepo
}

func (i *InMemoryDataStore) GetKeyRepository() KeyRepositoryInterface {
	return i.keyRepo
}

func (i *InMemoryDataStore) GetUserAccountRepository() UserAccountRepositoryInterface {
	return i.userAccountRepo
}
