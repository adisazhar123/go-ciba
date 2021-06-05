package service

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/adisazhar123/go-ciba/domain"
	"github.com/adisazhar123/go-ciba/grant"
	"github.com/adisazhar123/go-ciba/repository"
	"github.com/adisazhar123/go-ciba/service/http_auth"
	"github.com/adisazhar123/go-ciba/service/transport"
	"github.com/adisazhar123/go-ciba/util"
)

const (
	logTag = "[GO-CIBA CIBA SERVICE]"
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

	ValidateUserCode       func(code, givenCode string) bool
	ValidateBindingMessage func(bindingMessage string) bool
}

func defaultValidateUserCode(code, givenCode string) bool {
	return code == givenCode
}

func defaultValidateBindingMessage(bindingMessage string) bool {
	return bindingMessage != "" && len(bindingMessage) < 10
}

func NewAuthenticationRequest(r *http.Request) *AuthenticationRequest {
	authRequest := &AuthenticationRequest{
		ValidateUserCode:       defaultValidateUserCode,
		ValidateBindingMessage: defaultValidateBindingMessage,
	}
	_ = r.ParseForm()
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

func (ar *AuthenticationRequest) SetValidateUserCodeFunction(fn func(code, givenCode string) bool) *AuthenticationRequest {
	ar.ValidateUserCode = fn
	return ar
}

func (ar *AuthenticationRequest) SetValidateBindingMessageFunction(fn func(bindingMessage string) bool) *AuthenticationRequest {
	ar.ValidateBindingMessage = fn
	return ar
}

type AuthenticationResponse struct {
	AuthReqId string `json:"auth_req_id"`
	ExpiresIn int    `json:"expires_in"`
	Interval  *int   `json:"interval,omitempty"`
}

func makeSuccessfulAuthenticationResponse(authReqId string, expiresIn int, interval *int) *AuthenticationResponse {
	return &AuthenticationResponse{
		AuthReqId: authReqId,
		ExpiresIn: expiresIn,
		Interval:  interval,
	}
}

type ConsentRequest struct {
	AuthReqId string `json:"auth_req_id"`
	Consented *bool
}

type CibaServiceInterface interface {
	GrantServiceInterface
	ConsentServiceInterface
}

type CibaService struct {
	clientAppRepo   repository.ClientApplicationRepositoryInterface
	userAccountRepo repository.UserAccountRepositoryInterface
	cibaSessionRepo repository.CibaSessionRepositoryInterface
	keyRepo         repository.KeyRepositoryInterface

	scopeUtil             util.ScopeUtil
	authenticationContext *http_auth.ClientAuthenticationContext

	clientApp *domain.ClientApplication
	grant     grant.GrantTypeInterface

	notificationClient transport.NotificationInterface

	clientAppNotification transport.NotificationInterface

	tokenConfig *TokenConfig

	validateClientNotificationToken func(token string) bool

	mutex sync.Mutex
}

func defaultValidateClientNotificationToken(token string) bool {
	return token != ""
}

func (cs *CibaService) HandleAuthenticationRequest(request *AuthenticationRequest) (*AuthenticationResponse, *util.OidcError) {
	err := cs.ValidateAuthenticationRequestParameters(request)
	if err != nil {
		return nil, err
	}

	// Create new ciba session
	ciba := domain.NewCibaSession(cs.clientApp, request.LoginHint, request.BindingMessage, request.ClientNotificationToken, request.Scope, request.RequestedExpiry, cs.grant.(*grant.CibaGrant).PollInterval)
	if err := cs.cibaSessionRepo.Create(ciba); err != nil {
		log.Println("An error occurred", err)
		return nil, util.ErrGeneral
	}

	firebaseToken := "dPI3eiikuIOyVtqdfVsSzf:APA91bGuu6JA9ROzHXuywum7HWLTmxNmbz9y45Ma50q26sEJZ4T7LoPvqm8aiuah_MM2WhQPmbIo3h2o0FBzlJ6n1VuQdu_HXHlXexy2eUjtwuWXdWyFz3wjUeCsR7Rvn_3jqrVyp5F-"

	if err := cs.notificationClient.Send(map[string]interface{}{
		"to":               firebaseToken,
		"data.auth_req_id": ciba.AuthReqId,
	}); err != nil {
		log.Printf("[go-ciba][cibaservice] an error occured sending consent to user %s", err.Error())
		return nil, util.ErrGeneral
	}

	return makeSuccessfulAuthenticationResponse(ciba.AuthReqId, ciba.ExpiresIn, ciba.Interval), nil
}

func (cs *CibaService) ValidateAuthenticationRequestParameters(request *AuthenticationRequest) *util.OidcError {
	// Make sure client application exists
	clientApp, err := cs.clientAppRepo.FindById(request.ClientId)
	if err != nil {
		return util.ErrGeneral
	}
	if clientApp == nil {
		return util.ErrUnauthorizedClient
	}
	cs.clientApp = clientApp

	// Make sure authentication type is correct e.g. http_auth basic, client secret JWT etc.
	clientAuth := cs.authenticationContext.AuthenticateClient(request.r, clientApp)
	if !clientAuth {
		return util.ErrInvalidClient
	}
	// Make sure client app is registered to use CIBA
	if !clientApp.IsRegisteredToUseGrantType(grant.IdentifierCiba) {
		return util.ErrUnauthorizedClient
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
		return util.ErrInvalidRequest
	}

	// Make sure hint is valid, it must correspond to a valid user
	user, err := cs.userAccountRepo.FindById(request.LoginHint)
	if err != nil {
		return util.ErrGeneral
	}
	if user == nil {
		return util.ErrUnknownUserId
	}

	// Make sure scope is valid for chosen client
	if !cs.scopeUtil.ScopeExist(clientApp.GetScope(), request.Scope) {
		return util.ErrInvalidScope
	}

	// Client registered using ping or push must provide client_notification_token
	if (clientApp.GetTokenMode() == domain.ModePing || clientApp.GetTokenMode() == domain.ModePush) && !cs.validateClientNotificationToken(request.ClientNotificationToken) {
		log.Printf("%s client notification is missing or not well formed\n", logTag)
		return util.ErrInvalidRequest
	}

	if !request.ValidateBindingMessage(request.BindingMessage) {
		return util.ErrInvalidBindingMessage
	}

	// Client registered using user code must supply user_code
	if clientApp.IsUserCodeSupported() && request.UserCode == "" {
		return util.ErrMissingUserCode
	}

	// Check if user code is correct
	if clientApp.IsUserCodeSupported() && !request.ValidateUserCode(user.GetUseCode(), request.UserCode) {
		return util.ErrInvalidUserCode
	}

	return nil
}

func (cs *CibaService) HandleConsentRequest(request *ConsentRequest) error {
	cibaSession, err := cs.cibaSessionRepo.FindById(request.AuthReqId)

	if err != nil {
		// not valid
		return util.ErrGeneral
	}
	if cibaSession == nil {
		return errors.New("ciba session not found")
	}

	clientApp, err := cs.clientAppRepo.FindById(cibaSession.ClientId)
	if err != nil {
		return err
	}
	if clientApp == nil {
		return errors.New("client app not found")
	}

	if !cibaSession.Valid || cibaSession.Consented != nil || cibaSession.IsTimeExpired() {
		// not valid
		log.Printf("[go-ciba][cibaservice] ciba session %s isn't valid\n", cibaSession.AuthReqId)
		if clientApp.TokenMode == domain.ModePush {
			_ = cs.clientAppNotification.Send(map[string]interface{}{
				"token_method":              domain.ModePush,
				"success":                   false,
				"oidc_error":                util.ErrExpiredToken,
				"endpoint":                  clientApp.ClientNotificationEndpoint,
				"client_notification_token": cibaSession.ClientNotificationToken,
			})
		}
		return util.ErrExpiredToken
	}
	cibaSession.Consented = request.Consented
	if err := cs.cibaSessionRepo.Update(cibaSession); err != nil {
		// not valid
		return err
	}

	if request.Consented != nil && *request.Consented && clientApp.TokenMode == domain.ModePush {
		extraClaims := make(map[string]interface{})
		now := util.NowInt()

		key, err := cs.keyRepo.FindPrivateKeyByClientId(cibaSession.ClientId)

		if err != nil {
			return err
		}

		if key == nil {
			log.Printf("%s cannot find key for client ID %s", logTag, cibaSession.ClientId)
			return util.ErrInvalidGrant
		}

		extraClaims["urn:openid:params:jwt:claim:auth_req_id"] = cibaSession.AuthReqId

		tokens := cs.grant.(*grant.CibaGrant).CreateAccessTokenAndIdToken(domain.DefaultCibaIdTokenClaims{
			DefaultIdTokenClaims: domain.DefaultIdTokenClaims{
				Aud:      cibaSession.ClientId,
				AuthTime: now,
				Iat:      now,
				Exp:      cs.tokenConfig.IdTokenLifeTime,
				Iss:      cs.tokenConfig.Issuer,
				Sub:      cibaSession.UserId,
			},
			AuthReqId: cibaSession.AuthReqId,
		}, extraClaims, key.Private, key.Alg, key.ID)

		cibaSession.Expire()
		cibaSession.IdToken = tokens.IdToken.Value

		if err := cs.cibaSessionRepo.Update(cibaSession); err != nil {
			log.Printf("[go-ciba][pushtoken] failed updating CIBA session. %s", err.Error())
			return util.ErrGeneral
		}

		_ = cs.clientAppNotification.Send(map[string]interface{}{
			"token_method":              domain.ModePush,
			"success":                   true,
			"auth_req_id":               cibaSession.AuthReqId,
			"access_token":              tokens.AccessToken.Value,
			"token_type":                tokens.AccessToken.TokenType,
			"expires_in":                tokens.AccessToken.ExpiresIn,
			"id_token":                  tokens.IdToken.Value,
			"client_notification_token": cibaSession.ClientNotificationToken,
			"endpoint":                  clientApp.ClientNotificationEndpoint,
		})
	} else if request.Consented != nil && !*request.Consented && clientApp.TokenMode == domain.ModePush {
		cibaSession.Consented = request.Consented

		if err := cs.cibaSessionRepo.Update(cibaSession); err != nil {
			log.Printf("[go-ciba][pushtoken] failed updating CIBA session. %s", err.Error())
			return util.ErrGeneral
		}

		_ = cs.clientAppNotification.Send(map[string]interface{}{
			"token_method":              domain.ModePush,
			"success":                   false,
			"oidc_error":                util.ErrAccessDenied,
			"client_notification_token": cibaSession.ClientNotificationToken,
			"endpoint":                  clientApp.ClientNotificationEndpoint,
		})
	} else if request.Consented != nil && clientApp.TokenMode == domain.ModePing {
		cibaSession.Consented = request.Consented

		if err := cs.cibaSessionRepo.Update(cibaSession); err != nil {
			log.Printf("[go-ciba][pushtoken] failed updating CIBA session. %s", err.Error())
			return util.ErrGeneral
		}

		_ = cs.clientAppNotification.Send(map[string]interface{}{
			"token_method":              domain.ModePing,
			"client_notification_token": cibaSession.ClientNotificationToken,
			"endpoint":                  clientApp.ClientNotificationEndpoint,
			"auth_req_id":               cibaSession.AuthReqId,
		})
	}

	return nil
}

func (cs *CibaService) GetGrantIdentifier() string {
	return cs.grant.GetIdentifier()
}
