package go_ciba

import (
	"fmt"
	"log"
	"net/http"

	"github.com/adisazhar123/go-ciba/grant"
	"github.com/adisazhar123/go-ciba/repository"
	"github.com/adisazhar123/go-ciba/service"
	"github.com/adisazhar123/go-ciba/util"
)

type AuthorizationServerInterface interface {
	AddGrant(grant grant.GrantTypeInterface)
	AddService(grantService service.GrantServiceInterface)
	HandleCibaRequest(r http.Request) (*service.AuthenticationResponse, *util.OidcError)
}

type AuthorizationServer struct {
	grantServices map[string]service.GrantServiceInterface
	dataStore     repository.DataStoreInterface
}

func NewAuthorizationServer(ds repository.DataStoreInterface) *AuthorizationServer {
	return &AuthorizationServer{
		grantServices: make(map[string]service.GrantServiceInterface),
		dataStore:     ds,
	}
}

func (as *AuthorizationServer) AddService(gs service.GrantServiceInterface) {
	_, exist := as.grantServices[gs.GetGrantIdentifier()]
	if !exist {
		as.grantServices[gs.GetGrantIdentifier()] = gs
		log.Printf("added grant type: %s\n", gs.GetGrantIdentifier())
	}
}

func (as *AuthorizationServer) HandleCibaRequest(request *service.AuthenticationRequest) (*service.AuthenticationResponse, *util.OidcError) {
	if _, exist := as.grantServices[grant.IdentifierCiba]; !exist {
		panic(fmt.Sprintf("grant %s doesn't exist", grant.IdentifierCiba))
	}
	return as.grantServices[grant.IdentifierCiba].HandleAuthenticationRequest(request)
}

func (as *AuthorizationServer) initCibaService() *service.CibaService {
	if as.dataStore == nil {
		panic("dataStore is null")
	}
	return &service.CibaService{}
}
