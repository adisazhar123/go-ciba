package service

import (
	"github.com/adisazhar123/ciba-server/domain"
	"github.com/adisazhar123/ciba-server/repository"
	"github.com/adisazhar123/ciba-server/util"
	"github.com/cockroachdb/errors"
	"net/http"
	"strconv"
)

type AuthenticationRequest struct {
	ClientId string
	ClientSecret string

	AcrValues string
	BindingMessage string
	ClientNotificationToken string
	IdTokenHint string
	LoginHint string
	LoginHintToken string
	RequestedExpiry int
	Scope string
	UserCode string
}

func MakeAuthenticationRequest(r *http.Request) *AuthenticationRequest {
	authRequest := &AuthenticationRequest{}
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

	return authRequest
}

type CibaServiceInterface interface {
	ValidateAuthenticationRequestParameters(request *AuthenticationRequest) error
}

type CibaService struct {
	clientAppRepo repository.ClientApplicationRepositoryInterface
	userAccountRepo repository.UserAccountRepositoryInterface
	scopeUtil util.ScopeUtil
}

func (cs *CibaService) ValidateAuthenticationRequestParameters(request *AuthenticationRequest) (interface{}, error) {
	// Make sure client application exists
	clientApp := cs.clientAppRepo.FindById(request.ClientId)
	if clientApp == nil {
		// TODO: return error client app not found
		return util.ErrUnauthorizedClient, errors.New(util.ErrUnauthorizedClient.ErrorDescription)
	}

	// Make authentication type is correct e.g. http basic, client secret JWT etc.

	// Make sure client app is registered to use CIBA

	// Validate JWT if request is signed

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
		return
	}

	// Make sure hint is valid, it must correspond to a valid user
	user, err := cs.userAccountRepo.FindById(request.LoginHint)
	if err != nil {
		panic("error userAccountRepo.FindById")
	}
	if user == nil {
		// TODO: return error user not found
		return util.ErrUnknownUserId, errors.New(util.ErrUnknownUserId.ErrorDescription)
	}

	// Make sure scope is valid for chosen client
	if !cs.scopeUtil.ScopeExist(clientApp, request.Scope) {
		// TODO: return scope invalid
		return util.ErrInvalidScope, errors.New(util.ErrInvalidScope.ErrorDescription)
	}

	// Client registered using ping or push must provide client_notification_token
	if clientApp.GetTokenMode() == domain.MODE_PING || clientApp.GetTokenMode() == domain.MODE_PUSH && request.ClientNotificationToken == "" {
		// TODO: return error client notification given
		return util.ErrInvalidRequest, errors.New(util.ErrInvalidRequest.ErrorDescription)
	}

	// TODO: Allow custom logic for binding message
	if request.BindingMessage != "" && len(request.BindingMessage) > 10 {
		// TODO: return error invalid binding message
		return util.ErrInvalidBindingMessage, errors.New(util.ErrInvalidBindingMessage.ErrorDescription)
	}

	// Client registered using user code must supply user_code
	if clientApp.IsUserCodeSupported() && request.UserCode == "" {
		// TODO: return error missing user code
		return util.ErrMissingUserCode, errors.New(util.ErrMissingUserCode.ErrorDescription)
	}

	// Check if user code is correct
	if clientApp.IsUserCodeSupported() && user.GetUseCode() != request.UserCode {
		// TODO: return error incorrect user code
		return util.ErrInvalidUserCode, errors.New(util.ErrInvalidUserCode.ErrorDescription)
	}
}