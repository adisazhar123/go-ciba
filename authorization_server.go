package go_ciba

import (
	"fmt"
	"github.com/adisazhar123/go-ciba/grant"
	"github.com/adisazhar123/go-ciba/service"
	"log"
	"net/http"
)

type AuthorizationServerInterface interface {
	AddGrant(grant grant.GrantTypeInterface)
	AddService(grantService service.GrantServiceInterface)
	HandleCibaRequest(r http.Request)
}

type AuthorizationServer struct {
	grantServices map[string]service.GrantServiceInterface
}

func (as *AuthorizationServer) AddService(gs service.GrantServiceInterface) {
	_, exist := as.grantServices[gs.GetGrantIdentifier()]
	if !exist {
		as.grantServices[gs.GetGrantIdentifier()] = gs
		log.Printf("added grant type: %s\n", gs.GetGrantIdentifier())
	}
}

func (as *AuthorizationServer) HandleCibaRequest(request *service.AuthenticationRequest) (interface{}, error) {
	if _, exist := as.grantServices[grant.IdentifierCiba]; !exist {
		panic(fmt.Sprintf("grant %s doesn't exist", grant.IdentifierCiba))
	}
	return as.grantServices[grant.IdentifierCiba].HandleAuthenticationRequest(request)
}

func (as *AuthorizationServer) ValidateAuthenticationRequest(request *service.AuthenticationRequest) {

}
