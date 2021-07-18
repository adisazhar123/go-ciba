package go_ciba

import (
	"log"

	"github.com/adisazhar123/go-ciba/grant"
	"github.com/adisazhar123/go-ciba/repository"
	"github.com/adisazhar123/go-ciba/service"
	"github.com/adisazhar123/go-ciba/util"
)

type AuthorizationServerInterface interface {
	AddGrant(grant grant.GrantTypeInterface)
	AddService(grantService service.GrantServiceInterface)
	HandleCibaRequest(request *service.AuthenticationRequest) (*service.AuthenticationResponse, *util.OidcError)
	HandleConsentRequest(request *service.ConsentRequest) *util.OidcError
}

type authorizationServer struct {
	grantServices map[string]service.GrantServiceInterface
	dataStore     repository.DataStoreInterface
}

func NewAuthorizationServer(ds repository.DataStoreInterface) *authorizationServer {
	return &authorizationServer{
		grantServices: make(map[string]service.GrantServiceInterface),
		dataStore:     ds,
	}
}

func (as *authorizationServer) AddService(gs service.GrantServiceInterface) {
	_, exist := as.grantServices[gs.GetGrantIdentifier()]
	if !exist {
		as.grantServices[gs.GetGrantIdentifier()] = gs
		log.Printf("added grant type: %s\n", gs.GetGrantIdentifier())
	}
}

func (as *authorizationServer) HandleCibaRequest(request *service.AuthenticationRequest) (*service.AuthenticationResponse, *util.OidcError) {
	if _, exist := as.grantServices[grant.IdentifierCiba]; !exist {
		return nil, util.ErrGeneral
	}
	cs, ok := as.grantServices[grant.IdentifierCiba].(service.CibaServiceInterface)
	if !ok {
		return nil, util.ErrGeneral
	}
	return cs.HandleAuthenticationRequest(request)
}

func (as *authorizationServer) HandleConsentRequest(request *service.ConsentRequest) *util.OidcError {
	if _, exist := as.grantServices[grant.IdentifierCiba]; !exist {
		return util.ErrGeneral
	}
	cs, ok := as.grantServices[grant.IdentifierCiba].(service.CibaServiceInterface)
	if !ok {
		return util.ErrGeneral
	}
	return cs.HandleConsentRequest(request)
}
