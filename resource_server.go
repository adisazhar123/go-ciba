package go_ciba

import (
	"net/http"
	"strings"

	"github.com/adisazhar123/go-ciba/repository"
	"github.com/adisazhar123/go-ciba/util"
)

type ResourceRequest struct {
	accessToken string
}

func getTokenFromHeader(h http.Header) string {
	// Implementation assumes that the access token
	// is stored as:
	// Authorization: Bearer access_token_here
	val := h.Get("Authorization")
	vals := strings.Split(val, " ")
	if len(vals) != 2 {
		return ""
	}
	return vals[1]
}

func NewResourceRequest(r *http.Request) *ResourceRequest {
	return &ResourceRequest{
		accessToken: getTokenFromHeader(r.Header),
	}
}

type ResourceServerInterface interface {
	HandleResourceRequest(r *ResourceRequest) error
}

type resourceServer struct {
	accessTokenRepo repository.AccessTokenRepositoryInterface
	scopeUtil       util.ScopeUtil
}

func (rs *resourceServer) HandleResourceRequest(r *ResourceRequest, scope string) *util.OidcError {
	token, err := rs.accessTokenRepo.Find(r.accessToken)
	if err != nil {
		return util.ErrGeneral
	}
	if token == nil {
		return util.ErrInvalidToken
	}
	if scope != "" && !rs.scopeUtil.ScopeExist(token.Scope, scope) {
		return util.ErrInsufficientScope
	}
	if token.IsExpired() {
		return util.ErrInvalidToken
	}
	return nil
}
