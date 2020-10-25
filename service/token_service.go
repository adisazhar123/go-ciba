 package service

 import (
	 "github.com/cockroachdb/errors"
	 "github.com/adisazhar123/go-ciba/repository"
	 "github.com/adisazhar123/go-ciba/util"
	 "log"
	 "net/http"
 )

 type TokenRequest struct {
	 clientId string
	 clientSecret string
	 grantType string
	 authReqId string
	 httpMethod string
 }

 func NewTokenRequest(r *http.Request) *TokenRequest {
	 tokenRequest := &TokenRequest{}
	 r.ParseForm()
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
 	clientAppRepo repository.ClientApplicationRepositoryInterface
 	cibaSessionRepo repository.CibaSessionRepositoryInterface
 }

 func (t *TokenService) HandleTokenRequest(request *TokenRequest) (interface{}, error) {
	 panic("implement me")
 }

 func (t *TokenService) GrantAccessToken(request *TokenRequest) (interface{}, error) {
	 cs, err := t.cibaSessionRepo.FindById(request.authReqId)
	 if err != nil {
	 	log.Fatalln(err)
	 	return util.ErrGeneral, errors.New(util.ErrGeneral.ErrorDescription)
	 }
	 if cs == nil {
	 	return false, err
	 }
 }

