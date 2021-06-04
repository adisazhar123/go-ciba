package service

import (
	"log"
	"net/http"
	"time"

	"github.com/adisazhar123/go-ciba/domain"
	"github.com/adisazhar123/go-ciba/grant"
	"github.com/adisazhar123/go-ciba/repository"
	"github.com/adisazhar123/go-ciba/util"
)

type TokenRequest struct {
	clientId     string
	clientSecret string
	grantType    string
	authReqId    string
	httpMethod   string
}

const (
	LogTag = "[GO-CIBA TOKEN SERVICE]"
)

func NewTokenRequest(r *http.Request) *TokenRequest {
	tokenRequest := &TokenRequest{}
	_ = r.ParseForm()
	form := r.Form

	tokenRequest.authReqId = form.Get("auth_req_id")
	tokenRequest.grantType = form.Get("grant_type")
	tokenRequest.httpMethod = r.Method

	return tokenRequest
}

type TokenServiceInterface interface {
	HandleTokenRequest(request *TokenRequest) (interface{}, *util.OidcError)
	ValidateTokenRequest(request *TokenRequest) (interface{}, *util.OidcError)
	GrantAccessToken(request *TokenRequest) (interface{}, *util.OidcError)
}

type TokenConfig struct {
	PollingInterval     int
	AccessTokenLifeTime int
	IdTokenLifeTime     int
	Issuer              string
}

type TokenService struct {
	accessTokenRepo repository.AccessTokenRepositoryInterface
	clientAppRepo   repository.ClientApplicationRepositoryInterface
	cibaSessionRepo repository.CibaSessionRepositoryInterface
	keyRepo         repository.KeyRepositoryInterface
	config          *TokenConfig
	// TODO: support other grant types as well, not just CIBA.
	grant *grant.CibaGrant
}

func (t *TokenService) HandleTokenRequest(request *TokenRequest) (interface{}, *util.OidcError) {
	return t.GrantAccessToken(request)
}

func (t *TokenService) validate(cs *domain.CibaSession) (interface{}, *util.OidcError) {
	if !cs.IsValid() || cs.IsTimeExpired() {
		return util.ErrExpiredToken, util.ErrExpiredToken
	} else if cs.IsAuthorizationPending() {
		return util.ErrAuthorizationPending, util.ErrAuthorizationPending
	} else if !cs.IsConsented() {
		return util.ErrAccessDenied, util.ErrAccessDenied
	}
	return true, nil
}

type UserConsentResponse struct {
	err    *util.OidcError
	status bool
}

func waitForUserConsent(response chan UserConsentResponse, authReqId string, cibaSessionRepo repository.CibaSessionRepositoryInterface) {
	start := int(time.Now().Unix())
	timeout := 30
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
			now := int(time.Now().Unix())
			timeTaken := now - start
			if timeTaken > timeout {
				log.Printf("%s waiting for user consent hit timeout\n", LogTag)
				response <- UserConsentResponse{
					err:    util.ErrAuthorizationPending,
					status: false,
				}
				break
			}
			time.Sleep(1)
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

func (t *TokenService) GrantAccessToken(request *TokenRequest) (*domain.Tokens, *util.OidcError) {
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
		now := int(time.Now().Unix())
		// This CIBA session has requested a token before - not the first time.
		if cs.LatestTokenRequestedAt != nil {
			reqInterval := now - *cs.LatestTokenRequestedAt

			// Make sure that the time between the last token request
			// and the current token request isn't too quick
			if reqInterval < t.config.PollingInterval {
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
		_, err := t.validate(cs)
		if err != nil {
			return nil, err
		}
	} else if ca.TokenMode == domain.ModePush {
		return nil, util.ErrUnauthorizedClient
	}

	key, err := t.keyRepo.FindPrivateKeyByClientId(request.clientId)

	if key == nil {
		log.Printf("%s cannot find key for client ID %s", LogTag, request.clientId)
		return nil, util.ErrInvalidGrant
	}
	// TODO: support extra claims
	extraClaims := make(map[string]interface{})
	now := util.NowInt()
	// TODO: support other grant types as well, not just CIBA.
	tokens := t.grant.CreateAccessTokenAndIdToken(domain.DefaultCibaIdTokenClaims{
		DefaultIdTokenClaims: domain.DefaultIdTokenClaims{
			Aud:      request.clientId,
			AuthTime: now,
			Iat:      now,
			Exp:      t.config.IdTokenLifeTime,
			Iss:      t.config.Issuer,
			Sub:      cs.UserId,
		},
		AuthReqId: request.authReqId,
	}, extraClaims, key.Private, key.Alg, key.ID)
	// value, clientId, userId, scope string, expires int
	accessToken := domain.NewAccessToken(tokens.AccessToken.Value, request.clientId, cs.UserId, cs.Scope, now+tokens.AccessToken.ExpiresIn)
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
