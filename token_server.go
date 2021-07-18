package go_ciba

import (
	"github.com/adisazhar123/go-ciba/service"
	"github.com/adisazhar123/go-ciba/util"
)

type TokenServerInterface interface {
	HandleTokenRequest(request *service.TokenRequest) (*service.TokenResponse, *util.OidcError)
}

type tokenServer struct {
	service service.TokenServiceInterface
}

func (t *tokenServer) HandleTokenRequest(request *service.TokenRequest) (*service.TokenResponse, *util.OidcError) {
	return t.service.HandleTokenRequest(request)
}

