package service

import (
	"github.com/adisazhar123/go-ciba/repository"
	"github.com/adisazhar123/go-ciba/util"
	"github.com/cockroachdb/errors"
	"log"
	"net/http"
)

type TokenRequest struct {
	clientId     string
	clientSecret string
	grantType    string
	authReqId    string
	httpMethod   string
}

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
	GrantAccessToken(request *TokenRequest) (interface{}, error)
}

type TokenService struct {
	clientAppRepo   repository.ClientApplicationRepositoryInterface
	cibaSessionRepo repository.CibaSessionRepositoryInterface
}

func (t *TokenService) HandleTokenRequest(request *TokenRequest) (interface{}, error) {
	panic("implement me")
}

func (t *TokenService) GrantAccessToken(request *TokenRequest) (interface{}, error) {
	// Do some validation
	// Check if auth_req_id exists
	cs, err := t.cibaSessionRepo.FindById(request.authReqId)
	if err != nil {
		log.Fatalln(err)
		return util.ErrGeneral, err
	}
	if cs == nil {
		return false, err
	}
	// Check if client_id that is attached to auth_req_id is registered to use CIBA
	ca, err := t.clientAppRepo.FindById(cs.ClientId)
	if err != nil {
		log.Fatalln(err)
		return util.ErrGeneral, err
	}
	if ca == nil {
		return util.ErrInvalidClient, errors.New(util.ErrInvalidClient.ErrorDescription)
	}

	if !cs.IsValid() || cs.IsTimeExpired() {
		return util.ErrExpiredToken, errors.New(util.ErrExpiredToken.ErrorDescription)
	} else if cs.IsAuthorizationPending() {
		return util.ErrAuthorizationPending, errors.New(util.ErrAuthorizationPending.ErrorDescription)
	} else if !cs.IsConsented() {
		return util.ErrAccessDenied, errors.New(util.ErrAccessDenied.ErrorDescription)
	}


}
