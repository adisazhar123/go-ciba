package ciba_server

import (
	"github.com/adisazhar123/ciba-server/grant"
	"github.com/adisazhar123/ciba-server/service"
	"log"
	"net/http"
)

type AuthorizationServerInterface interface {
	AddGrant(grant grant.GrantTypeInterface)
	HandleCibaRequest(r http.Request)
}

type AuthorizationServer struct {
	grants map[string]grant.GrantTypeInterface
	cibaService service.CibaServiceInterface
}

func (as *AuthorizationServer) AddGrant(gt grant.GrantTypeInterface) {
	_, exist := as.grants[gt.GetIdentifier()]
	if !exist {
		as.grants[gt.GetIdentifier()] = gt
		log.Printf("added grant type: %s\n", gt.GetIdentifier())
	}
}

func (as *AuthorizationServer) HandleCibaRequest(request *service.AuthenticationRequest) {
	if err := as.grants[grant.IDENTIFIER_CIBA].ValidateAuthenticationRequest(request); err != nil {
		panic(err)
	}
}

func (as *AuthorizationServer) ValidateAuthenticationRequest(request *service.AuthenticationRequest) {

}