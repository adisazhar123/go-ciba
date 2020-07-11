package service

import (
	"github.com/adisazhar123/ciba-server/domain"
	"github.com/adisazhar123/ciba-server/repository"
	"github.com/adisazhar123/ciba-server/util"
)

type AuthenticationRequest struct {
	ClientId string
	ClientSecret string

	ClientNotificationToken string
	AcrValues string
	BindingMessage string
	UserCode string
	RequestedExpiry int
	Scope string
	LoginHintToken string
	IdTokenHint string
	LoginHint string
}

type CibaService struct {
	clientAppRepo repository.ClientApplicationRepositoryInterface
	userAccountRepo repository.UserAccountRepositoryInterface
	scopeUtil util.ScopeUtil
}

func (cs *CibaService) validateAuthenticationRequestParameters(request AuthenticationRequest) {
	// Make sure client application exists
	clientApp := cs.clientAppRepo.FindById(request.ClientId)
	if clientApp == nil {
		// TODO: return error client app not found
		return
	}

	// authentication type is correct e.g. http basic, client secret JWT etc.

	// make sure client app is registered to use CIBA

	// validate JWT if request is signed


	// validate all authentication request parameters
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
		return
	}

	// Make sure scope is valid for chosen client
	if !cs.scopeUtil.ScopeExist(clientApp, request.Scope) {
		// TODO: return scope invalid
		return
	}

	// Client registered using ping or push must provide client_notification_token
	if clientApp.GetTokenMode() == domain.MODE_PING || clientApp.GetTokenMode() == domain.MODE_PUSH && request.ClientNotificationToken == "" {
		// TODO: return error client notification given
		return
	}

	// TODO: Allow custom logic for binding message
	if request.BindingMessage != "" && len(request.BindingMessage) > 10 {
		// TODO: return error invalid binding message
		return
	}

	// Client registered using user code must supply user_code
	if clientApp.IsUserCodeSupported() && request.UserCode == "" {
		// TODO: return error missing user code
		return
	}

	// Check if user code is correct
	if clientApp.IsUserCodeSupported() && user.GetUseCode() != request.UserCode {
		// TODO: return error incorrect user code
		return
	}
}