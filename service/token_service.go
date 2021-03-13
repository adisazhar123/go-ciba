package service

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/adisazhar123/go-ciba/domain"
	"github.com/adisazhar123/go-ciba/grant"
	"github.com/adisazhar123/go-ciba/repository"
	"github.com/adisazhar123/go-ciba/util"
	"github.com/cockroachdb/errors"
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
	HandleTokenRequest(request *TokenRequest) (interface{}, error)
	ValidateTokenRequest(request *TokenRequest) (interface{}, error)
	GrantAccessToken(request *TokenRequest) (interface{}, error)
}

type TokenConfig struct {
	PollingInterval     int
	Alg                 string
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

func (t *TokenService) HandleTokenRequest(request *TokenRequest) (interface{}, error) {
	return t.GrantAccessToken(request)
}

func (t *TokenService) validate(cs *domain.CibaSession) (interface{}, error) {
	if !cs.IsValid() || cs.IsTimeExpired() {
		return util.ErrExpiredToken, errors.New(util.ErrExpiredToken.ErrorDescription)
	} else if cs.IsAuthorizationPending() {
		return util.ErrAuthorizationPending, errors.New(util.ErrAuthorizationPending.ErrorDescription)
	} else if !cs.IsConsented() {
		return util.ErrAccessDenied, errors.New(util.ErrAccessDenied.ErrorDescription)
	}
	return true, nil
}

type UserConsentResponse struct {
	err     error
	errType interface{}
	status  bool
}

func waitForUserConsent(response chan UserConsentResponse, authReqId string, cibaSessionRepo repository.CibaSessionRepositoryInterface) {
	start := int(time.Now().Unix())
	timeout := 30
	for true {
		cs, err := cibaSessionRepo.FindById(authReqId)
		if err != nil {
			log.Println(err)
			response <- UserConsentResponse{
				err:     err,
				errType: nil,
				status:  false,
			}
			break
		}
		if cs.IsAuthorizationPending() {
			now := int(time.Now().Unix())
			if (start - now) > timeout {
				log.Printf("%s Waiting for user consent hit timeout\n", LogTag)
				response <- UserConsentResponse{
					err:     errors.New(util.ErrAuthorizationPending.ErrorDescription),
					errType: util.ErrAuthorizationPending,
					status:  false,
				}
				break
			}
			time.Sleep(1)
			continue
		} else if cs.IsConsented() {
			log.Printf("%s User has consented\n", LogTag)
			response <- UserConsentResponse{
				err:     nil,
				errType: nil,
				status:  true,
			}
			break
		} else {
			log.Printf("%s User didn't give consent\n", LogTag)
			response <- UserConsentResponse{
				err:     errors.New(util.ErrAccessDenied.ErrorDescription),
				errType: util.ErrAccessDenied,
				status:  false,
			}
			break
		}
	}
}

func (t *TokenService) GrantAccessToken(request *TokenRequest) (interface{}, error) {
	// Do some validation
	// Check if auth_req_id exists
	cs, err := t.cibaSessionRepo.FindById(request.authReqId)
	if err != nil {
		log.Println(err)
		return util.ErrGeneral, err
	}
	if cs == nil {
		return false, err
	}
	// Check if client_id that is attached to auth_req_id is registered to use CIBA
	ca, err := t.clientAppRepo.FindById(cs.ClientId)
	if err != nil {
		log.Println(err)
		return util.ErrGeneral, err
	}
	if ca == nil {
		return util.ErrInvalidClient, errors.New(util.ErrInvalidClient.ErrorDescription)
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
				return util.ErrSlowDown, errors.New(util.ErrSlowDown.ErrorDescription)
			}
		}

		cs.LatestTokenRequestedAt = &now
		if err := t.cibaSessionRepo.Update(cs); err != nil {
			log.Printf("%s Failed updating CIBA session.", LogTag)
			return util.ErrGeneral, err
		}

		ucrChan := make(chan UserConsentResponse)
		go waitForUserConsent(ucrChan, request.authReqId, t.cibaSessionRepo)
		resp := <-ucrChan
		if resp.err != nil {
			log.Printf("%s Failed waiting for user consent. %s", LogTag, resp.err.Error())
			return resp.errType, resp.err
		}
	} else if ca.TokenMode == domain.ModePing {
		res, err := t.validate(cs)
		if err != nil {
			return res, err
		}
	} else if ca.TokenMode == domain.ModePush {
		return util.ErrUnauthorizedClient, errors.New(util.ErrUnauthorizedClient.ErrorDescription)
	}

	key, err := t.keyRepo.FindPrivateKeyByClientId(request.clientId)

	if key == nil {
		log.Printf("%s Cannot find key for client ID. %s", LogTag, request.clientId)
		return util.ErrGeneral, fmt.Errorf("cannot find key for client ID %s", request.clientId)
	}

	extraClaims := make(map[string]interface{})
	now := int(time.Now().Unix())
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
	}, extraClaims, key.Private, t.config.Alg, key.ID)
	// value, clientId, userId, scope string, expires int
	accessToken := domain.NewAccessToken(tokens.AccessToken.Value, request.clientId, cs.UserId, cs.Scope, now+tokens.AccessToken.ExpiresIn)
	if err := t.accessTokenRepo.Create(accessToken); err != nil {
		log.Printf("%s Cannot create access token. %s", LogTag, err.Error())
		return util.ErrGeneral, err
	}

	cs.Expire()
	cs.IdToken = tokens.IdToken.Value
	if err := t.cibaSessionRepo.Update(cs); err != nil {
		log.Printf("%s Failed updating CIBA session. %s", LogTag, err.Error())
		return util.ErrGeneral, err
	}

	return tokens, nil
}
