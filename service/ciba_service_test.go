package service

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"testing"

	"github.com/adisazhar123/go-ciba/domain"
	"github.com/adisazhar123/go-ciba/grant"
	"github.com/adisazhar123/go-ciba/service/http_auth"
	"github.com/adisazhar123/go-ciba/test_data"
	"github.com/adisazhar123/go-ciba/util"
	"github.com/stretchr/testify/assert"
)

type UserAccountVolatileRepository struct {
	data map[string]*domain.UserAccount
}

// In memory mock of UserAccountRepositoryInterface.
func newUserAccountVolatileRepository() *UserAccountVolatileRepository {
	return &UserAccountVolatileRepository{
		data: map[string]*domain.UserAccount{
			fmt.Sprintf("user_account:%s", test_data.User1.Id): &test_data.User1,
			fmt.Sprintf("user_account:%s", test_data.User2.Id): &test_data.User2,
			fmt.Sprintf("user_account:%s", test_data.User3.Id): &test_data.User3,
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

// Create new ClientAuthenticationContext object.
func newAuthenticationContext() *http_auth.ClientAuthenticationContext {
	return &http_auth.ClientAuthenticationContext{}
}

type notificationClientMock struct { }

func (n notificationClientMock) Send(data map[string]interface{}) error {
	return nil
}

func newCibaService() *CibaService {
	return &CibaService{
		clientAppRepo:         test_data.NewClientApplicationVolatileRepository(),
		userAccountRepo:       newUserAccountVolatileRepository(),
		cibaSessionRepo:       test_data.NewCibaSessionVolatileRepository(),
		scopeUtil:             util.ScopeUtil{},
		authenticationContext: newAuthenticationContext(),
		grant:                 grant.NewCibaGrant(),
		notificationClient: &notificationClientMock{},
	}
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
	cs := newCibaService()
	auth := createAuthorizationHeaderBasic(test_data.ClientAppPing.Id, test_data.ClientAppPing.Secret)

	form := url.Values{}
	form.Set("scope", test_data.ClientAppPing.Scope)
	form.Set("client_notification_token", util.GenerateRandomString())
	form.Set("login_hint", test_data.User1.Id)
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
	assert.NotEmpty(t, authRes.AuthReqId)
}

// Tests a Ciba request with client application registered as ping mode and also requires a user code
// user code parameter is given
// expected to succeed/ no error.
func TestCibaService_HandleAuthenticationRequest_Valid_WithUserCode_Ping(t *testing.T) {
	cs := newCibaService()
	auth := createAuthorizationHeaderBasic(test_data.ClientAppPingUserCodeSupported.Id, test_data.ClientAppPingUserCodeSupported.Secret)

	form := url.Values{}
	form.Set("scope", test_data.ClientAppPingUserCodeSupported.Scope)
	form.Set("client_notification_token", util.GenerateRandomString())
	form.Set("login_hint", test_data.User3.Id)
	form.Set("binding_message", "aa-123")
	form.Set("requested_expiry", "120")
	form.Set("user_code", "1999")

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
	assert.NotEmpty(t, authRes.AuthReqId)
}

// TODO: Test for poll mode

// Tests a Ciba request with client application registered as ping mode
// the authorization head built has incorrect client_id so the authentication
// will fail. Expected return error unauthorized_client.
func TestCibaService_HandleAuthenticationRequest_Invalid_ClientId_ClientCredentials_Ping(t *testing.T) {
	cs := newCibaService()
	auth := createAuthorizationHeaderBasic(test_data.ClientAppPing.Id+"break-id", test_data.ClientAppPing.Secret)

	form := url.Values{}
	form.Set("scope", test_data.ClientAppPing.Scope)
	form.Set("client_notification_token", util.GenerateRandomString())
	form.Set("login_hint", test_data.User1.Id)
	form.Set("binding_message", "aa-123")
	form.Set("requested_expiry", "120")

	request, _ := http.NewRequest(http.MethodPost, "ciba.example.com/bc-authorize", strings.NewReader(form.Encode()))
	request.Header.Add("Authorization", fmt.Sprintf("Basic %s", auth))
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Add("Content-Length", strconv.Itoa(len(form.Encode())))

	authReq := NewAuthenticationRequest(request)
	res, err := cs.HandleAuthenticationRequest(authReq)

	assert.EqualError(t, err, util.ErrUnauthorizedClient.Error())
	assert.Nil(t, res)
}

// Tests a Ciba request with client application registered as ping mode
// the authorization head built has incorrect client_secret so the authentication
// will fail invalid_client.
func TestCibaService_HandleAuthenticationRequest_Invalid_Password_ClientCredentials_Ping(t *testing.T) {
	cs := newCibaService()
	auth := createAuthorizationHeaderBasic(test_data.ClientAppPing.Id, test_data.ClientAppPing.Secret+"break-secret")

	form := url.Values{}
	form.Set("scope", test_data.ClientAppPing.Scope)
	form.Set("client_notification_token", util.GenerateRandomString())
	form.Set("login_hint", test_data.User1.Id)
	form.Set("binding_message", "aa-123")
	form.Set("requested_expiry", "120")

	request, _ := http.NewRequest(http.MethodPost, "ciba.example.com/bc-authorize", strings.NewReader(form.Encode()))
	request.Header.Add("Authorization", fmt.Sprintf("Basic %s", auth))
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Add("Content-Length", strconv.Itoa(len(form.Encode())))

	authReq := NewAuthenticationRequest(request)
	res, err := cs.HandleAuthenticationRequest(authReq)

	assert.Nil(t, res)
	assert.EqualError(t, err, util.ErrInvalidClient.Error())
}

// Tests a Ciba request with client application registered as ping mode
// multiple login hints are used: login_hint, login_hint_token, id_token_hint
// expected return error invalid_request.
func TestCibaService_HandleAuthenticationRequest_Invalid_MultipleHints_Ping(t *testing.T) {
	cs := newCibaService()
	auth := createAuthorizationHeaderBasic(test_data.ClientAppPing.Id, test_data.ClientAppPing.Secret)

	form := url.Values{}
	form.Set("scope", test_data.ClientAppPing.Scope)
	form.Set("client_notification_token", util.GenerateRandomString())
	form.Set("login_hint", test_data.User1.Id)
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

	assert.Nil(t, res)
	assert.EqualError(t, err, util.ErrInvalidRequest.Error())
}

// Tests a Ciba request with client application registered as ping mode
// client is requesting a scope that isn't registered
// expected return error invalid_scope.
func TestCibaService_HandleAuthenticationRequest_Invalid_UnregisteredScope_Ping(t *testing.T) {
	cs := newCibaService()
	auth := createAuthorizationHeaderBasic(test_data.ClientAppPing.Id, test_data.ClientAppPing.Secret)

	form := url.Values{}
	form.Set("scope", "openid tree random")
	form.Set("client_notification_token", util.GenerateRandomString())
	form.Set("login_hint", test_data.User1.Id)
	form.Set("binding_message", "aa-123")
	form.Set("requested_expiry", "120")

	request, _ := http.NewRequest(http.MethodPost, "ciba.example.com/bc-authorize", strings.NewReader(form.Encode()))
	request.Header.Add("Authorization", fmt.Sprintf("Basic %s", auth))
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Add("Content-Length", strconv.Itoa(len(form.Encode())))

	authReq := NewAuthenticationRequest(request)
	res, err := cs.HandleAuthenticationRequest(authReq)

	assert.Nil(t, res)
	assert.EqualError(t, err, util.ErrInvalidScope.Error())
}

// Tests a Ciba request with client application not registered to use Ciba
// this client is registered to use authorization code and client credentials
// expected return error unauthorized_client
func TestCibaService_HandleAuthenticationRequest_Invalid_ClientAppNotRegisteredToUseCiba(t *testing.T) {
	cs := newCibaService()
	auth := createAuthorizationHeaderBasic(test_data.ClientAppNotRegisteredToUseCiba.Id, test_data.ClientAppNotRegisteredToUseCiba.Secret)

	form := url.Values{}
	form.Set("scope", test_data.ClientAppNotRegisteredToUseCiba.Scope)
	form.Set("login_hint", test_data.User1.Id)
	form.Set("binding_message", "aa-123")
	form.Set("requested_expiry", "120")

	request, _ := http.NewRequest(http.MethodPost, "ciba.example.com/bc-authorize", strings.NewReader(form.Encode()))
	request.Header.Add("Authorization", fmt.Sprintf("Basic %s", auth))
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Add("Content-Length", strconv.Itoa(len(form.Encode())))

	authReq := NewAuthenticationRequest(request)
	res, err := cs.HandleAuthenticationRequest(authReq)

	assert.Nil(t, res)
	assert.EqualError(t, err, util.ErrUnauthorizedClient.Error())
}

// Tests a Ciba request with a client app registered as ping mode
// this request is done without a client notification parameter,
// which is required for ping and push modes
// expected return error invalid_request
func TestCibaService_HandleAuthenticationRequest_Invalid_ClientAppPingWithoutNotificationToken(t *testing.T) {
	cs := newCibaService()
	auth := createAuthorizationHeaderBasic(test_data.ClientAppPing.Id, test_data.ClientAppPing.Secret)

	form := url.Values{}
	form.Set("scope", test_data.ClientAppPing.Scope)
	form.Set("login_hint", test_data.User1.Id)
	form.Set("binding_message", "aa-123")
	form.Set("requested_expiry", "120")

	request, _ := http.NewRequest(http.MethodPost, "ciba.example.com/bc-authorize", strings.NewReader(form.Encode()))
	request.Header.Add("Authorization", fmt.Sprintf("Basic %s", auth))
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Add("Content-Length", strconv.Itoa(len(form.Encode())))

	authReq := NewAuthenticationRequest(request)
	res, err := cs.HandleAuthenticationRequest(authReq)

	assert.Nil(t, res)
	assert.EqualError(t, err, util.ErrInvalidRequest.Error())
}

// Tests a Ciba request with a client app registered as push mode
// this request is done without a client notification parameter,
// which is required for ping and push modes
// expected return error invalid_request
func TestCibaService_HandleAuthenticationRequest_Invalid_ClientAppPushWithoutNotificationToken(t *testing.T) {
	cs := newCibaService()
	auth := createAuthorizationHeaderBasic(test_data.ClientAppPush.Id, test_data.ClientAppPush.Secret)

	form := url.Values{}
	form.Set("scope", test_data.ClientAppPush.Scope)
	form.Set("login_hint", test_data.User1.Id)
	form.Set("binding_message", "aa-123")
	form.Set("requested_expiry", "120")

	request, _ := http.NewRequest(http.MethodPost, "ciba.example.com/bc-authorize", strings.NewReader(form.Encode()))
	request.Header.Add("Authorization", fmt.Sprintf("Basic %s", auth))
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Add("Content-Length", strconv.Itoa(len(form.Encode())))

	authReq := NewAuthenticationRequest(request)
	res, err := cs.HandleAuthenticationRequest(authReq)

	assert.Nil(t, res)
	assert.EqualError(t, err, util.ErrInvalidRequest.Error())
}

// Tests a Ciba request with a client app registered as push mode
// this request includes a binding message with length > 10
// a binding message should be concise
// expected return error invalid_binding_message
func TestCibaService_HandleAuthenticationRequest_Invalid_BindingMessageTooLong(t *testing.T) {
	cs := newCibaService()
	auth := createAuthorizationHeaderBasic(test_data.ClientAppPush.Id, test_data.ClientAppPush.Secret)

	form := url.Values{}
	form.Set("scope", test_data.ClientAppPush.Scope)
	form.Set("login_hint", test_data.User1.Id)
	form.Set("binding_message", "aa-123-321-123-321-123-321-123")
	form.Set("requested_expiry", "120")
	form.Set("client_notification_token", "41217fd5-10dc-46e8-8151-27b7edf372fa")

	request, _ := http.NewRequest(http.MethodPost, "ciba.example.com/bc-authorize", strings.NewReader(form.Encode()))
	request.Header.Add("Authorization", fmt.Sprintf("Basic %s", auth))
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Add("Content-Length", strconv.Itoa(len(form.Encode())))

	authReq := NewAuthenticationRequest(request)
	res, err := cs.HandleAuthenticationRequest(authReq)

	assert.Nil(t, res)
	assert.EqualError(t, err, util.ErrInvalidBindingMessage.Error())
}

func TestCibaService_HandleAuthenticationRequest_Invalid_UserCodeNotGiven(t *testing.T) {
	cs := newCibaService()
	auth := createAuthorizationHeaderBasic(test_data.ClientAppPushUserCodeSupported.Id, test_data.ClientAppPushUserCodeSupported.Secret)

	form := url.Values{}
	form.Set("scope", test_data.ClientAppPushUserCodeSupported.Scope)
	form.Set("login_hint", test_data.User1.Id)
	form.Set("binding_message", "aa-123")
	form.Set("requested_expiry", "120")
	form.Set("client_notification_token", "41217fd5-10dc-46e8-8151-27b7edf372fa")

	request, _ := http.NewRequest(http.MethodPost, "ciba.example.com/bc-authorize", strings.NewReader(form.Encode()))
	request.Header.Add("Authorization", fmt.Sprintf("Basic %s", auth))
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Add("Content-Length", strconv.Itoa(len(form.Encode())))

	authReq := NewAuthenticationRequest(request)
	res, err := cs.HandleAuthenticationRequest(authReq)

	assert.Nil(t, res)
	assert.EqualError(t, err, util.ErrMissingUserCode.Error())
}

func TestCibaService_HandleAuthenticationRequest_Invalid_WrongUserCode(t *testing.T) {
	cs := newCibaService()
	auth := createAuthorizationHeaderBasic(test_data.ClientAppPushUserCodeSupported.Id, test_data.ClientAppPushUserCodeSupported.Secret)

	form := url.Values{}
	form.Set("scope", test_data.ClientAppPushUserCodeSupported.Scope)
	form.Set("login_hint", test_data.User3.Id)
	form.Set("binding_message", "aa-123")
	form.Set("requested_expiry", "120")
	form.Set("client_notification_token", "41217fd5-10dc-46e8-8151-27b7edf372fa")
	form.Set("user_code", "1234")

	request, _ := http.NewRequest(http.MethodPost, "ciba.example.com/bc-authorize", strings.NewReader(form.Encode()))
	request.Header.Add("Authorization", fmt.Sprintf("Basic %s", auth))
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Add("Content-Length", strconv.Itoa(len(form.Encode())))

	authReq := NewAuthenticationRequest(request)
	res, err := cs.HandleAuthenticationRequest(authReq)

	assert.Nil(t, res)
	assert.EqualError(t, err, util.ErrInvalidUserCode.Error())
}

func TestCibaService_ValidateAuthenticationRequestParameters(t *testing.T) {

}
