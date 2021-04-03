package test_data

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/adisazhar123/go-ciba/domain"
	"github.com/adisazhar123/go-ciba/grant"
	"github.com/adisazhar123/go-ciba/service/http_auth"
)

type ClientApplicationVolatileRepository struct {
	data map[string]*domain.ClientApplication
}

var (
	privateKeyFile, _ = os.Open("../test_data/key.pem")
	publicKeyFile, _  = os.Open("../test_data/public.pem")

	privateKey, _ = ioutil.ReadAll(privateKeyFile)
	publicKey, _  = ioutil.ReadAll(publicKeyFile)

	Key1 = domain.Key{
		ID:       "1",
		ClientId: "unknown",
		Alg:      "RS256",
		Public:   string(publicKey),
		Private:  string(privateKey),
	}

	Key2 = domain.Key{
		ID:       "2",
		ClientId: CibaSession9.ClientId,
		Alg:      "RS256",
		Public:   string(publicKey),
		Private:  string(privateKey),
	}

	Key3 = domain.Key{
		ID:       "3",
		ClientId: CibaSession10.ClientId,
		Alg:      "RS256",
		Public:   string(publicKey),
		Private:  string(privateKey),
	}

	Key4 = domain.Key{
		ID:       "4",
		ClientId: CibaSession11.ClientId,
		Alg:      "RS256",
		Public:   string(publicKey),
		Private:  string(privateKey),
	}

	Key5 = domain.Key{
		ID:       "5",
		ClientId: CibaSession12.ClientId,
		Alg:      "RS256",
		Public:   string(publicKey),
		Private:  string(privateKey),
	}

	Key6 = domain.Key{
		ID:       "6",
		ClientId: CibaSession13.ClientId,
		Alg:      "RS256",
		Public:   string(publicKey),
		Private:  string(privateKey),
	}

	// Client applications
	// non signed, non user code
	ClientAppPush = domain.ClientApplication{
		Id:                              "8df692eb-968c-4ba0-8a7c-c082d5a56982",
		Secret:                          "secret",
		Name:                            "client-app-push",
		Scope:                           "openid email profile",
		TokenMode:                       domain.ModePush,
		ClientNotificationEndpoint:      "go-ciba.dev/notification",
		AuthenticationRequestSigningAlg: "",
		UserCodeParameterSupported:      false,
		TokenEndpointAuthMethod:         http_auth.ClientSecretBasic,
		GrantTypes:                      fmt.Sprintf("%s", grant.IdentifierCiba),
	}

	ClientAppPing = domain.ClientApplication{
		Id:                              "420d637b-ff22-4e48-88fb-237aa2131e72",
		Secret:                          "secret",
		Name:                            "client-app-ping",
		Scope:                           "openid email profile",
		TokenMode:                       domain.ModePing,
		ClientNotificationEndpoint:      "go-ciba.dev/notification",
		AuthenticationRequestSigningAlg: "",
		UserCodeParameterSupported:      false,
		TokenEndpointAuthMethod:         http_auth.ClientSecretBasic,
		GrantTypes:                      fmt.Sprintf("%s", grant.IdentifierCiba),
	}

	ClientAppPoll = domain.ClientApplication{
		Id:                         "f07aa98e-d072-4c0c-a71c-bb6d070fb002",
		Secret:                     "secret",
		Name:                       "client-app-poll",
		Scope:                      "openid email profile",
		TokenMode:                  domain.ModePoll,
		UserCodeParameterSupported: false,
		GrantTypes:                 fmt.Sprintf("%s", grant.IdentifierCiba),
	}

	ClientAppPushUserCodeSupported = domain.ClientApplication{
		Id:                              "e2d9bcd7-0f5a-47b6-8017-b50537e98330",
		Secret:                          "secret",
		Name:                            "client-app-push-user-code",
		Scope:                           "openid email profile",
		TokenMode:                       domain.ModePush,
		ClientNotificationEndpoint:      "go-ciba.dev/notification",
		AuthenticationRequestSigningAlg: "",
		UserCodeParameterSupported:      true,
		TokenEndpointAuthMethod:         http_auth.ClientSecretBasic,
		GrantTypes:                      fmt.Sprintf("%s", grant.IdentifierCiba),
	}

	ClientAppPingUserCodeSupported = domain.ClientApplication{
		Id:                              "5dd6f0fc-75a2-4dee-873e-a55eceb0c3ee",
		Secret:                          "secret",
		Name:                            "client-app-ping-user-code",
		Scope:                           "openid email profile",
		TokenMode:                       domain.ModePing,
		ClientNotificationEndpoint:      "go-ciba.dev/notification",
		AuthenticationRequestSigningAlg: "",
		UserCodeParameterSupported:      true,
		TokenEndpointAuthMethod:         http_auth.ClientSecretBasic,
		GrantTypes:                      fmt.Sprintf("%s", grant.IdentifierCiba),
	}

	// not registered to use ciba
	ClientAppNotRegisteredToUseCiba = domain.ClientApplication{
		Id:                              "aa27b00d-04ba-4021-97b0-eacf8b013126",
		Secret:                          "secret",
		Name:                            "client-app-auth-code",
		Scope:                           "openid email profile",
		TokenMode:                       "",
		ClientNotificationEndpoint:      "",
		AuthenticationRequestSigningAlg: "",
		UserCodeParameterSupported:      false,
		RedirectUri:                     "clientapp.dev/redirect",
		TokenEndpointAuthMethod:         http_auth.ClientSecretBasic,
		TokenEndpointAuthSigningAlg:     "",
		GrantTypes:                      "authorization_code client_credentials",
		PublicKeyUri:                    "",
	}

	// Users
	User1 = domain.UserAccount{
		Id:        "59f37eab-39a6-4e87-9dd4-2a29194f09a4",
		Name:      "user-1",
		Email:     "user-1@email.com",
		Password:  "secret",
		UserCode:  "",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	User2 = domain.UserAccount{
		Id:        "b4e6ba16-d09c-46b3-9feb-96e4f2e396f3",
		Name:      "user-2",
		Email:     "user-2@email.com",
		Password:  "secret",
		UserCode:  "",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	User3 = domain.UserAccount{
		Id:        "ba714f46-a3c1-496f-8267-1da563472d4d",
		Name:      "user-3",
		Email:     "user-3@email.com",
		Password:  "secret",
		UserCode:  "1999",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	consent     = true
	notConsent  = false
	expiresLong = 9999999

	CibaSession1 = domain.CibaSession{
		AuthReqId: "8f8080f3-7b3e-4f86-af93-60fc14392008",
		ClientId:  "e334bdee-3c48-4c98-b96b-2e5251ca7ad4",
	}

	CibaSession2 = domain.CibaSession{
		AuthReqId: "aba53bfb-9ac3-4f0b-8191-f019668b0601",
		ClientId:  "unknown-client-id",
	}

	CibaSession3 = domain.CibaSession{
		AuthReqId: "af1005ed-4b7a-475e-ad99-66ba397ea70f",
		ClientId:  ClientAppNotRegisteredToUseCiba.Id,
	}

	CibaSession4 = domain.CibaSession{
		AuthReqId: "ad2221df-d758-435f-a17e-c70fe0db00b6",
		ClientId:  ClientAppPush.Id,
	}

	CibaSession5 = domain.CibaSession{
		AuthReqId: "0fd7c8fd-4e27-479a-9503-0e7f2fff2b0b",
		ClientId:  ClientAppPing.Id,
		ExpiresIn: 3600,
		Valid:     true,
		CreatedAt: time.Now().Add(-1 * time.Duration(5) * time.Hour),
	}

	CibaSession6 = domain.CibaSession{
		AuthReqId: "8c125df0-9079-42ca-9da8-279d6c75335d",
		ClientId:  ClientAppPing.Id,
		ExpiresIn: expiresLong,
		Valid:     true,
		Consented: nil,
		CreatedAt: time.Now(),
	}

	CibaSession7 = domain.CibaSession{
		AuthReqId: "43b9afb7-c1f6-4017-a077-be3bb62563cb",
		ClientId:  ClientAppPing.Id,
		ExpiresIn: expiresLong,
		Valid:     true,
		Consented: &notConsent,
		CreatedAt: time.Now(),
	}

	CibaSession8 = domain.CibaSession{
		AuthReqId: "2d17bd7d-701c-40c0-bd92-1d85b50fa3ba",
		ClientId:  ClientAppPing.Id,
		ExpiresIn: expiresLong,
		Valid:     true,
		Consented: &consent,
		CreatedAt: time.Now(),
	}

	CibaSession9 = domain.CibaSession{
		AuthReqId: "9b18f36d-d294-464e-b5bf-585c72718faf",
		ClientId:  ClientAppPingUserCodeSupported.Id,
		ExpiresIn: expiresLong,
		Valid:     true,
		Consented: &consent,
		CreatedAt: time.Now(),
	}

	CibaSession10 = domain.CibaSession{
		AuthReqId:              "1385561f-f542-432d-9f08-669c766f5051",
		ClientId:               ClientAppPoll.Id,
		ExpiresIn:              expiresLong,
		Valid:                  true,
		Consented:              &consent,
		LatestTokenRequestedAt: nil,
		CreatedAt:              time.Time{},
	}

	now = int(time.Now().Unix())

	CibaSession11 = domain.CibaSession{
		AuthReqId:              "38f2aec8-4cfc-4982-9bd7-35e09cc60916",
		ClientId:               ClientAppPoll.Id,
		ExpiresIn:              expiresLong,
		Valid:                  true,
		Consented:              nil,
		LatestTokenRequestedAt: &now,
		CreatedAt:              time.Time{},
	}

	CibaSession12 = domain.CibaSession{
		AuthReqId:              "f0325001-569a-4cfd-8a0f-0338b0055064",
		ClientId:               ClientAppPoll.Id,
		ExpiresIn:              expiresLong,
		Valid:                  true,
		Consented:              nil,
		LatestTokenRequestedAt: nil,
		CreatedAt:              time.Time{},
	}

	CibaSession13 = domain.CibaSession{
		AuthReqId:              "aaf07d23-f414-4ba2-8ecc-5b13db3b36f2",
		ClientId:               ClientAppPoll.Id,
		ExpiresIn:              expiresLong,
		Valid:                  true,
		Consented:              &notConsent,
		LatestTokenRequestedAt: nil,
		CreatedAt:              time.Time{},
	}
)

// In memory mock of ClientApplicationRepositoryInterface.
func NewClientApplicationVolatileRepository() *ClientApplicationVolatileRepository {
	return &ClientApplicationVolatileRepository{
		data: map[string]*domain.ClientApplication{
			fmt.Sprintf("client_application:%s", ClientAppPush.Id):                   &ClientAppPush,
			fmt.Sprintf("client_application:%s", ClientAppPing.Id):                   &ClientAppPing,
			fmt.Sprintf("client_application:%s", ClientAppNotRegisteredToUseCiba.Id): &ClientAppNotRegisteredToUseCiba,
			fmt.Sprintf("client_application:%s", ClientAppPushUserCodeSupported.Id):  &ClientAppPushUserCodeSupported,
			fmt.Sprintf("client_application:%s", ClientAppPingUserCodeSupported.Id):  &ClientAppPingUserCodeSupported,
			fmt.Sprintf("client_application:%s", ClientAppPoll.Id):                   &ClientAppPoll,
		},
	}
}

func (c *ClientApplicationVolatileRepository) Register(clientApp *domain.ClientApplication) error {
	key := fmt.Sprintf("client_application:%s", clientApp.Id)
	c.data[key] = clientApp
	return nil
}

func (c *ClientApplicationVolatileRepository) FindById(id string) (*domain.ClientApplication, error) {
	key := fmt.Sprintf("client_application:%s", id)
	clientApp, _ := c.data[key]
	return clientApp, nil
}

type CibaSessionVolatileRepository struct {
	data map[string]*domain.CibaSession
}

// In memory mock of CibaSessionRepositoryInterface.
func NewCibaSessionVolatileRepository() *CibaSessionVolatileRepository {
	return &CibaSessionVolatileRepository{data: map[string]*domain.CibaSession{
		fmt.Sprintf("%s", CibaSession1.AuthReqId):  &CibaSession1,
		fmt.Sprintf("%s", CibaSession2.AuthReqId):  &CibaSession2,
		fmt.Sprintf("%s", CibaSession3.AuthReqId):  &CibaSession3,
		fmt.Sprintf("%s", CibaSession4.AuthReqId):  &CibaSession4,
		fmt.Sprintf("%s", CibaSession5.AuthReqId):  &CibaSession5,
		fmt.Sprintf("%s", CibaSession6.AuthReqId):  &CibaSession6,
		fmt.Sprintf("%s", CibaSession7.AuthReqId):  &CibaSession7,
		fmt.Sprintf("%s", CibaSession8.AuthReqId):  &CibaSession8,
		fmt.Sprintf("%s", CibaSession9.AuthReqId):  &CibaSession9,
		fmt.Sprintf("%s", CibaSession10.AuthReqId): &CibaSession10,
		fmt.Sprintf("%s", CibaSession11.AuthReqId): &CibaSession11,
		fmt.Sprintf("%s", CibaSession12.AuthReqId): &CibaSession12,
		fmt.Sprintf("%s", CibaSession13.AuthReqId): &CibaSession13,
	}}
}

func (c CibaSessionVolatileRepository) FindById(id string) (*domain.CibaSession, error) {
	return c.data[id], nil
}

func (c CibaSessionVolatileRepository) Update(cibaSession *domain.CibaSession) error {
	c.data[cibaSession.AuthReqId] = cibaSession
	return nil
}

func (c CibaSessionVolatileRepository) Create(cibaSession *domain.CibaSession) error {
	key := fmt.Sprintf("%s", cibaSession.AuthReqId)
	c.data[key] = cibaSession
	return nil
}

type KeyVolatileRepository struct {
	data map[string]*domain.Key
}

func NewKeyVolatileRepository() *KeyVolatileRepository {
	defer privateKeyFile.Close()
	defer publicKeyFile.Close()

	return &KeyVolatileRepository{data: map[string]*domain.Key{
		fmt.Sprintf("%s", Key1.ID): &Key1,
		fmt.Sprintf("%s", Key2.ID): &Key2,
		fmt.Sprintf("%s", Key3.ID): &Key3,
		fmt.Sprintf("%s", Key4.ID): &Key4,
		fmt.Sprintf("%s", Key5.ID): &Key5,
		fmt.Sprintf("%s", Key6.ID): &Key6,
	}}
}

func (k KeyVolatileRepository) FindPrivateKeyByClientId(clientId string) (*domain.Key, error) {
	for _, v := range k.data {
		if v.ClientId == clientId {
			return v, nil
		}
	}
	return nil, nil
}
