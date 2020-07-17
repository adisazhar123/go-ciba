package service

import (
	"github.com/adisazhar123/ciba-server/domain"
	"github.com/adisazhar123/ciba-server/grant"
	"github.com/adisazhar123/ciba-server/repository"
	"github.com/adisazhar123/ciba-server/service/http_auth"
	"github.com/adisazhar123/ciba-server/util"
	"github.com/cockroachdb/errors"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type AuthenticationRequest struct {
	ClientId     string
	ClientSecret string

	AcrValues               string
	BindingMessage          string
	ClientNotificationToken string
	IdTokenHint             string
	LoginHint               string
	LoginHintToken          string
	RequestedExpiry         int
	Scope                   string
	UserCode                string
	Interval                int

	request string // holds signed request content

	r *http.Request
}

func NewAuthenticationRequest(r *http.Request) *AuthenticationRequest {
	authRequest := &AuthenticationRequest{}
	r.ParseForm()
	form := r.Form

	authRequest.AcrValues = form.Get("acr_values")
	authRequest.BindingMessage = form.Get("binding_message")
	authRequest.ClientNotificationToken = form.Get("client_notification_token")
	authRequest.IdTokenHint = form.Get("id_token_hint")
	authRequest.LoginHint = form.Get("login_hint")
	authRequest.LoginHintToken = form.Get("login_hint_token")
	expiry, _ := strconv.Atoi(form.Get("requested_expiry"))
	authRequest.RequestedExpiry = expiry
	authRequest.Scope = form.Get("scope")
	authRequest.UserCode = form.Get("user_code")

	credentials := http_auth.UtilGetClientCredentials(r)

	authRequest.ClientId = credentials.ClientId
	authRequest.ClientSecret = credentials.ClientSecret

	authRequest.r = r

	return authRequest
}

type AuthenticationResponse struct {
	AuthReqId string `json:"auth_req_id"`
	ExpiresIn int    `json:"expires_in"`
	Interval  int    `json:"interval,omitempty"`
}

func makeSuccessfulAuthenticationResponse(authReqId string, expiresIn, interval int) *AuthenticationResponse {
	return &AuthenticationResponse{
		AuthReqId: authReqId,
		ExpiresIn: expiresIn,
		Interval:  interval,
	}
}

type CibaServiceInterface interface {
	GrantServiceInterface
}

type CibaService struct {
	clientAppRepo   repository.ClientApplicationRepositoryInterface
	userAccountRepo repository.UserAccountRepositoryInterface
	cibaSessionRepo repository.CibaSessionRepositoryInterface

	scopeUtil             util.ScopeUtil
	authenticationContext *http_auth.ClientAuthenticationContext

	clientApp *domain.ClientApplication
	grant     grant.GrantTypeInterface
}

func (cs *CibaService) HandleAuthenticationRequest(request *AuthenticationRequest) (interface{}, error) {
	validation, err := cs.ValidateAuthenticationRequestParameters(request)
	if err != nil {
		return validation, err
	}

	// Create new ciba session
	ciba := domain.NewCibaSession(request.LoginHint, request.BindingMessage, request.ClientNotificationToken, request.Scope, request.RequestedExpiry, request.Interval)
	if err := cs.cibaSessionRepo.Create(ciba); err != nil {
		log.Fatalln("An error occurred", err)
		return nil, err
	}
	// TODO: Send via notification for push mode
	return makeSuccessfulAuthenticationResponse(ciba.AuthReqId, ciba.ExpiresIn, ciba.Interval), nil
}

func (cs *CibaService) ValidateAuthenticationRequestParameters(request *AuthenticationRequest) (interface{}, error) {
	// Make sure client application exists
	clientApp := cs.clientAppRepo.FindById(request.ClientId)
	if clientApp == nil {
		return util.ErrUnauthorizedClient, errors.New(util.ErrUnauthorizedClient.ErrorDescription)
	}

	// Make sure authentication type is correct e.g. http_auth basic, client secret JWT etc.
	clientAuth := cs.authenticationContext.AuthenticateClient(request.r, clientApp)
	if !clientAuth {
		return util.ErrInvalidClient, errors.New(util.ErrInvalidClient.ErrorDescription)
	}
	// Make sure client app is registered to use CIBA
	if !util.SliceStringContains(strings.Split(clientApp.GetGrantTypes(), " "), grant.IdentifierCiba) {
		return util.ErrUnauthorizedClient, errors.New(util.ErrUnauthorizedClient.ErrorDescription)
	}

	// TODO: Validate JWT if request is signed

	// Validate all authentication request parameters
	hintCounter := 0
	if request.LoginHintToken != "" {
		hintCounter++
	}
	if request.LoginHint != "" {
		hintCounter++
	}
	if request.IdTokenHint != "" {
		hintCounter++
	}
	// Make sure only one type of hint
	if hintCounter == 0 || hintCounter > 1 {
		// TODO: return error login hint == 0 || hint > 1
		log.Println("login hint failed", request)
		return util.ErrInvalidRequest, errors.New(util.ErrInvalidRequest.ErrorDescription)
	}

	// Make sure hint is valid, it must correspond to a valid user
	user, err := cs.userAccountRepo.FindById(request.LoginHint)
	if err != nil {
		panic("error userAccountRepo.FindById")
	}
	if user == nil {
		return util.ErrUnknownUserId, errors.New(util.ErrUnknownUserId.ErrorDescription)
	}

	// Make sure scope is valid for chosen client
	if !cs.scopeUtil.ScopeExist(clientApp.GetScope(), request.Scope) {
		return util.ErrInvalidScope, errors.New(util.ErrInvalidScope.ErrorDescription)
	}

	// Client registered using ping or push must provide client_notification_token
	if (request.ClientNotificationToken == "") && (clientApp.GetTokenMode() == domain.ModePing || clientApp.GetTokenMode() == domain.ModePush) {
		log.Println("failed here notif")
		return util.ErrInvalidRequest, errors.New(util.ErrInvalidRequest.ErrorDescription)
	}

	// TODO: Allow custom logic for binding message
	if request.BindingMessage != "" && len(request.BindingMessage) > 10 {
		return util.ErrInvalidBindingMessage, errors.New(util.ErrInvalidBindingMessage.ErrorDescription)
	}

	// Client registered using user code must supply user_code
	if clientApp.IsUserCodeSupported() && request.UserCode == "" {
		return util.ErrMissingUserCode, errors.New(util.ErrMissingUserCode.ErrorDescription)
	}

	// Check if user code is correct
	if clientApp.IsUserCodeSupported() && user.GetUseCode() != request.UserCode {
		return util.ErrInvalidUserCode, errors.New(util.ErrInvalidUserCode.ErrorDescription)
	}

	return true, nil
}

func (cs *CibaService) GetGrantIdentifier() string {
	return cs.grant.GetIdentifier()
}
