package go_ciba

import (
	"net/http"
	"strings"

	"github.com/adisazhar123/go-ciba/repository"
	"github.com/adisazhar123/go-ciba/util"
)

type ResourceRequest struct {
	accessToken string
	scope *string
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

func NewResourceRequest(r *http.Request, scope *string) *ResourceRequest {
	return &ResourceRequest{
		accessToken: getTokenFromHeader(r.Header),
		scope: scope,
	}
}

type ResourceServerInterface interface {
	HandleResourceRequest(r *ResourceRequest) error
}

type ResourceServer struct {
	accessTokenRepo repository.AccessTokenRepositoryInterface
	scopeUtil util.ScopeUtil
}

func (rs *ResourceServer) HandleResourceRequest(r *ResourceRequest) *util.OidcError {
	token, err := rs.accessTokenRepo.Find(r.accessToken)
	if err != nil {
		return util.ErrGeneral
	}
	if token == nil {
		return util.ErrInvalidToken
	}
	if r.scope != nil && !rs.scopeUtil.ScopeExist(token.Scope, *r.scope) {
		return util.ErrInsufficientScope
	}
	if token.IsExpired() {
		return util.ErrInvalidToken
	}
	return nil
}