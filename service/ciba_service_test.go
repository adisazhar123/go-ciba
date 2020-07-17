package service

import (
	"encoding/base64"
	"fmt"
	"github.com/adisazhar123/ciba-server/domain"
	"github.com/adisazhar123/ciba-server/grant"
	"github.com/adisazhar123/ciba-server/service/http_auth"
	"github.com/adisazhar123/ciba-server/util"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"testing"
	"time"
)

type ClientApplicationVolatileRepository struct {
	data map[string]*domain.ClientApplication
}

var (
	// Client applications
	// non signed, non user code
	ClientAppPush = domain.ClientApplication{
		Id:                              "8df692eb-968c-4ba0-8a7c-c082d5a56982",
		Secret:                          "secret",
		Name:                            "client-app-push",
		Scope:                           "openid email profile",
		TokenMode:                       domain.ModePush,
		ClientNotificationEndpoint:      "ciba-server.dev/notification",
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
		ClientNotificationEndpoint:      "ciba-server.dev/notification",
		AuthenticationRequestSigningAlg: "",
		UserCodeParameterSupported:      false,
		TokenEndpointAuthMethod:         http_auth.ClientSecretBasic,
		GrantTypes:                      fmt.Sprintf("%s", grant.IdentifierCiba),
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
)

// In memory mock of ClientApplicationRepositoryInterface.
func newClientApplicationVolatileRepository() *ClientApplicationVolatileRepository {
	return &ClientApplicationVolatileRepository{
		data: map[string]*domain.ClientApplication{
			fmt.Sprintf("client_application:%s", ClientAppPush.Id): &ClientAppPush,
			fmt.Sprintf("client_application:%s", ClientAppPing.Id): &ClientAppPing,
		},
	}
}

func (c *ClientApplicationVolatileRepository) Register(clientApp *domain.ClientApplication) error {
	key := fmt.Sprintf("client_application:%s", clientApp.Id)
	c.data[key] = clientApp
	return nil
}

func (c *ClientApplicationVolatileRepository) FindById(id string) *domain.ClientApplication {
	key := fmt.Sprintf("client_application:%s", id)
	clientApp, exist := c.data[key]
	if !exist {
		return nil
	}
	return clientApp
}

type UserAccountVolatileRepository struct {
	data map[string]*domain.UserAccount
}

// In memory mock of UserAccountRepositoryInterface.
func newUserAccountVolatileRepository() *UserAccountVolatileRepository {
	return &UserAccountVolatileRepository{
		data: map[string]*domain.UserAccount{
			fmt.Sprintf("user_account:%s", User1.Id): &User1,
			fmt.Sprintf("user_account:%s", User2.Id): &User2,
			fmt.Sprintf("user_account:%s", User3.Id): &User3,
		},
	}
}

func (u UserAccountVolatileRepository) FindById(id string) (*domain.UserAccount, error) {
	key := fmt.Sprintf("user_account:%s", id)
	user, exist := u.data[key]
	if !exist {
		return nil, nil
	}
	return user, nil
}

type CibaSessionVolatileRepository struct {
	data map[string]*domain.CibaSession
}

/// In memory mock of CibaSessionRepositoryInterface.
func newCibaSessionVolatileRepository() *CibaSessionVolatileRepository {
	return &CibaSessionVolatileRepository{data: make(map[string]*domain.CibaSession)}
}

func (c CibaSessionVolatileRepository) Create(cibaSession *domain.CibaSession) error {
	key := fmt.Sprintf("ciba_session:%s", cibaSession.AuthReqId)
	c.data[key] = cibaSession
	return nil
}

// Create new ClientAuthenticationContext object.
func newAuthenticationContext() *http_auth.ClientAuthenticationContext {
	return &http_auth.ClientAuthenticationContext{}
}

// Util function to build authorization header basic.
func createAuthorizationHeaderBasic(id, password string) string {
	return base64.StdEncoding.EncodeToString([]byte(id + ":" + password))
}

// Make sure that Ciba identifier is correct.
func TestCibaService_GetGrantIdentifier(t *testing.T) {
	cs := &CibaService{
		grant: grant.NewCibaGrant(),
	}
	id := "urn:openid:params:grant-type:ciba"
	assert.Equal(t, id, cs.GetGrantIdentifier())
}

// Tests a Ciba request with client application registered as ping mode
// expected to succeed/ no error.
func TestCibaService_HandleAuthenticationRequest_Valid_Ping(t *testing.T) {
	cs := &CibaService{
		clientAppRepo:         newClientApplicationVolatileRepository(),
		userAccountRepo:       newUserAccountVolatileRepository(),
		cibaSessionRepo:       newCibaSessionVolatileRepository(),
		scopeUtil:             util.ScopeUtil{},
		authenticationContext: newAuthenticationContext(),
		clientApp:             nil,
		grant:                 grant.NewCibaGrant(),
	}
	auth := createAuthorizationHeaderBasic(ClientAppPing.Id, ClientAppPing.Secret)

	form := url.Values{}
	form.Set("scope", ClientAppPing.Scope)
	form.Set("client_notification_token", util.GenerateRandomString())
	form.Set("login_hint", User1.Id)
	form.Set("binding_message", "aa-123")
	form.Set("requested_expiry", "120")

	request, _ := http.NewRequest(http.MethodPost, "ciba.example.com/bc-authorize", strings.NewReader(form.Encode()))
	request.Header.Add("Authorization", fmt.Sprintf("Basic %s", auth))
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Add("Content-Length", strconv.Itoa(len(form.Encode())))

	authReq := NewAuthenticationRequest(request)
	res, err := cs.HandleAuthenticationRequest(authReq)

	authRes := res.(*AuthenticationResponse)

	assert.Empty(t, err)
	assert.Empty(t, authRes.Interval)
	assert.Equal(t, 120, authRes.ExpiresIn)
}

// Tests a Ciba request with client application registered as ping mode
// the authorization head built has incorrect client_id so the authentication
// will fail. Expected return error unauthorized_client.
func TestCibaService_HandleAuthenticationRequest_Invalid_ClientId_ClientCredentials_Ping(t *testing.T) {
	cs := &CibaService{
		clientAppRepo:         newClientApplicationVolatileRepository(),
		userAccountRepo:       newUserAccountVolatileRepository(),
		cibaSessionRepo:       newCibaSessionVolatileRepository(),
		scopeUtil:             util.ScopeUtil{},
		authenticationContext: newAuthenticationContext(),
		clientApp:             nil,
		grant:                 grant.NewCibaGrant(),
	}
	auth := createAuthorizationHeaderBasic(ClientAppPing.Id+"break-id", ClientAppPing.Secret)

	form := url.Values{}
	form.Set("scope", ClientAppPing.Scope)
	form.Set("client_notification_token", util.GenerateRandomString())
	form.Set("login_hint", User1.Id)
	form.Set("binding_message", "aa-123")
	form.Set("requested_expiry", "120")

	request, _ := http.NewRequest(http.MethodPost, "ciba.example.com/bc-authorize", strings.NewReader(form.Encode()))
	request.Header.Add("Authorization", fmt.Sprintf("Basic %s", auth))
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Add("Content-Length", strconv.Itoa(len(form.Encode())))

	authReq := NewAuthenticationRequest(request)
	res, err := cs.HandleAuthenticationRequest(authReq)

	authRes := res.(*util.OidcError)

	assert.NotNil(t, err)
	assert.Equal(t, util.ErrUnauthorizedClient.Error, authRes.Error)
	assert.Equal(t, util.ErrUnauthorizedClient.ErrorDescription, authRes.ErrorDescription)
	assert.Equal(t, util.ErrUnauthorizedClient.Code, authRes.Code)
}

// Tests a Ciba request with client application registered as ping mode
// the authorization head built has incorrect client_secret so the authentication
// will fail invalid_client.
func TestCibaService_HandleAuthenticationRequest_Invalid_Password_ClientCredentials_Ping(t *testing.T) {
	cs := &CibaService{
		clientAppRepo:         newClientApplicationVolatileRepository(),
		userAccountRepo:       newUserAccountVolatileRepository(),
		cibaSessionRepo:       newCibaSessionVolatileRepository(),
		scopeUtil:             util.ScopeUtil{},
		authenticationContext: newAuthenticationContext(),
		clientApp:             nil,
		grant:                 grant.NewCibaGrant(),
	}
	auth := createAuthorizationHeaderBasic(ClientAppPing.Id, ClientAppPing.Secret+"break-secret")

	form := url.Values{}
	form.Set("scope", ClientAppPing.Scope)
	form.Set("client_notification_token", util.GenerateRandomString())
	form.Set("login_hint", User1.Id)
	form.Set("binding_message", "aa-123")
	form.Set("requested_expiry", "120")

	request, _ := http.NewRequest(http.MethodPost, "ciba.example.com/bc-authorize", strings.NewReader(form.Encode()))
	request.Header.Add("Authorization", fmt.Sprintf("Basic %s", auth))
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Add("Content-Length", strconv.Itoa(len(form.Encode())))

	authReq := NewAuthenticationRequest(request)
	res, err := cs.HandleAuthenticationRequest(authReq)

	authRes := res.(*util.OidcError)

	assert.NotNil(t, err)
	assert.Equal(t, util.ErrInvalidClient.Error, authRes.Error)
	assert.Equal(t, util.ErrInvalidClient.ErrorDescription, authRes.ErrorDescription)
	assert.Equal(t, util.ErrInvalidClient.Code, authRes.Code)
}

// Tests a Ciba request with client application registered as ping mode
// multiple login hints are used: login_hint, login_hint_token, id_token_hint
// expected return error invalid_request.
func TestCibaService_HandleAuthenticationRequest_Invalid_MultipleHints_Ping(t *testing.T) {
	cs := &CibaService{
		clientAppRepo:         newClientApplicationVolatileRepository(),
		userAccountRepo:       newUserAccountVolatileRepository(),
		cibaSessionRepo:       newCibaSessionVolatileRepository(),
		scopeUtil:             util.ScopeUtil{},
		authenticationContext: newAuthenticationContext(),
		clientApp:             nil,
		grant:                 grant.NewCibaGrant(),
	}
	auth := createAuthorizationHeaderBasic(ClientAppPing.Id, ClientAppPing.Secret)

	form := url.Values{}
	form.Set("scope", ClientAppPing.Scope)
	form.Set("client_notification_token", util.GenerateRandomString())
	form.Set("login_hint", User1.Id)
	form.Set("login_hint_token", "some_token_bla_bla_bla")
	form.Set("id_token_hint", "dummy-long-id-token-63f515e1-7add-499c-9024-a9eb88b98711")
	form.Set("binding_message", "aa-123")
	form.Set("requested_expiry", "120")

	request, _ := http.NewRequest(http.MethodPost, "ciba.example.com/bc-authorize", strings.NewReader(form.Encode()))
	request.Header.Add("Authorization", fmt.Sprintf("Basic %s", auth))
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Add("Content-Length", strconv.Itoa(len(form.Encode())))

	authReq := NewAuthenticationRequest(request)
	res, err := cs.HandleAuthenticationRequest(authReq)

	authRes := res.(*util.OidcError)

	assert.NotNil(t, err)
	assert.Equal(t, util.ErrInvalidRequest.Error, authRes.Error)
	assert.Equal(t, util.ErrInvalidRequest.ErrorDescription, authRes.ErrorDescription)
	assert.Equal(t, util.ErrInvalidRequest.Code, authRes.Code)
}

// Tests a Ciba request with client application registered as ping mode
// multiple login hints are used: login_hint, login_hint_token, id_token_hint
// expected return error invalid_scope.
func TestCibaService_HandleAuthenticationRequest_Invalid_UnregisteredScope_Ping(t *testing.T) {
	cs := &CibaService{
		clientAppRepo:         newClientApplicationVolatileRepository(),
		userAccountRepo:       newUserAccountVolatileRepository(),
		cibaSessionRepo:       newCibaSessionVolatileRepository(),
		scopeUtil:             util.ScopeUtil{},
		authenticationContext: newAuthenticationContext(),
		clientApp:             nil,
		grant:                 grant.NewCibaGrant(),
	}
	auth := createAuthorizationHeaderBasic(ClientAppPing.Id, ClientAppPing.Secret)

	form := url.Values{}
	form.Set("scope", "openid tree random")
	form.Set("client_notification_token", util.GenerateRandomString())
	form.Set("login_hint", User1.Id)
	form.Set("binding_message", "aa-123")
	form.Set("requested_expiry", "120")

	request, _ := http.NewRequest(http.MethodPost, "ciba.example.com/bc-authorize", strings.NewReader(form.Encode()))
	request.Header.Add("Authorization", fmt.Sprintf("Basic %s", auth))
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Add("Content-Length", strconv.Itoa(len(form.Encode())))

	authReq := NewAuthenticationRequest(request)
	res, err := cs.HandleAuthenticationRequest(authReq)

	authRes := res.(*util.OidcError)

	assert.NotNil(t, err)
	assert.Equal(t, util.ErrInvalidScope.Error, authRes.Error)
	assert.Equal(t, util.ErrInvalidScope.ErrorDescription, authRes.ErrorDescription)
	assert.Equal(t, util.ErrInvalidScope.Code, authRes.Code)
}

func TestCibaService_ValidateAuthenticationRequestParameters(t *testing.T) {

}
