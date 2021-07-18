package service

import (
	"log"
	"net/http"
	"time"

	"github.com/adisazhar123/go-ciba/domain"
	"github.com/adisazhar123/go-ciba/grant"
	"github.com/adisazhar123/go-ciba/repository"
	"github.com/adisazhar123/go-ciba/service/http_auth"
	"github.com/adisazhar123/go-ciba/util"
)

type TokenRequest struct {
	clientId     string
	clientSecret string
	grantType    string
	authReqId    string
	httpMethod   string

	r *http.Request
}

const (
	LogTag           = "[GO-CIBA TOKEN SERVICE]"
	timeoutInSeconds = 30
)

func NewTokenRequest(r *http.Request) *TokenRequest {
	tokenRequest := &TokenRequest{}
	_ = r.ParseForm()
	form := r.Form

	tokenRequest.authReqId = form.Get("auth_req_id")
	tokenRequest.grantType = form.Get("grant_type")
	tokenRequest.httpMethod = r.Method
	tokenRequest.r = r

	http_auth.PopulateClientCredentials(r, &tokenRequest.clientId, &tokenRequest.clientSecret)

	return tokenRequest
}

type TokenServiceInterface interface {
	HandleTokenRequest(request *TokenRequest) (*TokenResponse, *util.OidcError)
	GrantAccessToken(request *TokenRequest) (*domain.Tokens, *util.OidcError)
	ValidateTokenRequest(request *TokenRequest) (*domain.Tokens, *util.OidcError)
}

type TokenResponse struct {
	AccessToken  string  `json:"access_token"`
	TokenType    string  `json:"token_type"`
	RefreshToken *string `json:"refresh_token"`
	ExpiresIn    int64   `json:"expires_in"`
	IdToken      string  `json:"id_token"`
}

type tokenService struct {
	accessTokenRepo repository.AccessTokenRepositoryInterface
	clientAppRepo   repository.ClientApplicationRepositoryInterface
	cibaSessionRepo repository.CibaSessionRepositoryInterface
	keyRepo         repository.KeyRepositoryInterface
	userClaimRepo repository.UserClaimRepositoryInterface
	// TODO: support other grant types as well, not just CIBA.
	grant                 *grant.CibaGrant
	authenticationContext *http_auth.ClientAuthenticationContext
}

func NewTokenService(accessTokenRepo repository.AccessTokenRepositoryInterface, clientAppRepo repository.ClientApplicationRepositoryInterface, cibaSessionRepo repository.CibaSessionRepositoryInterface, keyRepo repository.KeyRepositoryInterface, userClaimRepo repository.UserClaimRepositoryInterface, grant *grant.CibaGrant) *tokenService {
	return &tokenService{
		accessTokenRepo:       accessTokenRepo,
		clientAppRepo:         clientAppRepo,
		cibaSessionRepo:       cibaSessionRepo,
		keyRepo:               keyRepo,
		userClaimRepo: userClaimRepo,
		grant:                 grant,
		authenticationContext: http_auth.NewClientAuthenticationContext(grant.Config),
	}
}

func makeSuccessfulTokenResponse(tokens *domain.Tokens) *TokenResponse {
	return &TokenResponse{
		AccessToken:  tokens.AccessToken.Value,
		TokenType:    tokens.AccessToken.TokenType,
		RefreshToken: nil,
		ExpiresIn:    tokens.AccessToken.ExpiresIn,
		IdToken:      tokens.IdToken.Value,
	}
}

// This performs authentication on the client app
func (t *tokenService) ValidateTokenRequest(request *TokenRequest) *util.OidcError {
	if request.grantType != grant.IdentifierCiba {
		return util.ErrUnsupportedGrantType
	}
	ca, err := t.clientAppRepo.FindById(request.clientId)
	if err != nil {
		log.Println(err)
		return util.ErrGeneral
	} else if ca == nil {
		return util.ErrInvalidClient
	}

	ok := t.authenticationContext.AuthenticateClient(request.r, ca)
	if !ok {
		return util.ErrInvalidClient
	}
	return nil
}

func (t *tokenService) HandleTokenRequest(request *TokenRequest) (*TokenResponse, *util.OidcError) {
	if err := t.ValidateTokenRequest(request); err != nil {
		return nil, err
	}
	tokens, err := t.GrantAccessToken(request)
	if err != nil {
		return nil, err
	}
	return makeSuccessfulTokenResponse(tokens), nil
}

func (t *tokenService) validate(cs *domain.CibaSession) *util.OidcError {
	if !cs.IsValid() || cs.IsTimeExpired() {
		return util.ErrExpiredToken
	} else if cs.IsAuthorizationPending() {
		return util.ErrAuthorizationPending
	} else if !cs.IsConsented() {
		return util.ErrAccessDenied
	}
	return nil
}

type UserConsentResponse struct {
	err    *util.OidcError
	status bool
}

func waitForUserConsent(response chan UserConsentResponse, authReqId string, cibaSessionRepo repository.CibaSessionRepositoryInterface) {
	start := util.NowInt()
	for true {
		cs, err := cibaSessionRepo.FindById(authReqId)
		if err != nil {
			log.Println(err)
			response <- UserConsentResponse{
				err:    util.ErrGeneral,
				status: false,
			}
			break
		}
		if cs.IsAuthorizationPending() {
			now := util.NowInt()
			timeTakenInSeconds := now - start
			if timeTakenInSeconds > timeoutInSeconds {
				log.Printf("%s waiting for user consent hit timeoutInSeconds\n", LogTag)
				response <- UserConsentResponse{
					err:    util.ErrAuthorizationPending,
					status: false,
				}
				break
			}
			time.Sleep(1 * time.Second)
			continue
		} else if cs.IsConsented() {
			log.Printf("%s user has consented\n", LogTag)
			response <- UserConsentResponse{
				err:    nil,
				status: true,
			}
			break
		} else {
			log.Printf("%s user didn't give consent\n", LogTag)
			response <- UserConsentResponse{
				err:    util.ErrAccessDenied,
				status: false,
			}
			break
		}
	}
}

func (t *tokenService) GrantAccessToken(request *TokenRequest) (*domain.Tokens, *util.OidcError) {
	// Do some validation
	// Check if auth_req_id exists
	cs, err := t.cibaSessionRepo.FindById(request.authReqId)
	if err != nil {
		log.Println(err)
		return nil, util.ErrGeneral
	} else if cs == nil || (cs.ClientId != request.clientId) {
		return nil, util.ErrInvalidGrant
	}

	// Check if client_id that is attached to auth_req_id is registered to use CIBA
	ca, err := t.clientAppRepo.FindById(cs.ClientId)
	if err != nil {
		log.Println(err)
		return nil, util.ErrGeneral
	} else if ca == nil {
		return nil, util.ErrInvalidClient
	} else if !ca.IsRegisteredToUseGrantType(grant.IdentifierCiba) {
		return nil, util.ErrUnauthorizedClient
	}

	// POLL method is long polling
	if ca.TokenMode == domain.ModePoll {
		err := t.validate(cs)
		if err != nil && err != util.ErrAuthorizationPending {
			return nil, err
		}
		now := util.NowInt()
		// This CIBA session has requested a token before - not the first time.
		if cs.LatestTokenRequestedAt != nil {
			reqInterval := now - *cs.LatestTokenRequestedAt

			// Make sure that the time between the last token request
			// and the current token request isn't too quick
			if t.grant.Config.PollingIntervalInSeconds != nil && reqInterval < *t.grant.Config.PollingIntervalInSeconds {
				return nil, util.ErrSlowDown
			}
		}

		cs.LatestTokenRequestedAt = &now
		if err := t.cibaSessionRepo.Update(cs); err != nil {
			log.Printf("%s failed updating CIBA session.", LogTag)
			return nil, util.ErrGeneral
		}

		ucrChan := make(chan UserConsentResponse)
		go waitForUserConsent(ucrChan, request.authReqId, t.cibaSessionRepo)
		resp := <-ucrChan
		if resp.err != nil {
			log.Printf("%s failed waiting for user consent. %s", LogTag, resp.err.Error())
			return nil, resp.err
		}
	} else if ca.TokenMode == domain.ModePing {
		err := t.validate(cs)
		if err != nil {
			return nil, err
		}
	} else if ca.TokenMode == domain.ModePush {
		return nil, util.ErrUnauthorizedClient
	}

	key, err := t.keyRepo.FindPrivateKeyByClientId(request.clientId)

	if key == nil {
		log.Printf("%s cannot find key for client Id %s", LogTag, request.clientId)
		return nil, util.ErrInvalidGrant
	}

	extraClaims := t.userClaimRepo.GetUserClaims(cs.UserId, cs.Scope)
	now := util.NowInt()
	// TODO: support other grant types as well, not just CIBA.
	tokens := t.grant.CreateAccessTokenAndIdToken(domain.DefaultCibaIdTokenClaims{
		DefaultIdTokenClaims: domain.DefaultIdTokenClaims{
			Aud:      request.clientId,
			AuthTime: now,
			Iat:      now,
			Exp:      t.grant.Config.IdTokenLifetimeInSeconds,
			Iss:      t.grant.Config.Issuer,
			Sub:      cs.UserId,
		},
		AuthReqId: request.authReqId,
	}, extraClaims, key.Private, key.Alg, key.Id)

	accessToken := domain.NewAccessToken(tokens.AccessToken.Value, request.clientId, cs.Hint, cs.Scope, time.Unix(now+tokens.AccessToken.ExpiresIn, 0))
	if err := t.accessTokenRepo.Create(accessToken); err != nil {
		log.Printf("%s cannot create access token. %s", LogTag, err.Error())
		return nil, util.ErrGeneral
	}

	cs.Expire()
	cs.IdToken = tokens.IdToken.Value
	if err := t.cibaSessionRepo.Update(cs); err != nil {
		log.Printf("%s failed updating CIBA session. %s", LogTag, err.Error())
		return nil, util.ErrGeneral
	}

	return tokens, nil
}
