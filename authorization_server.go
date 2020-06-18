package ciba_server

import (
	"github.com/adisazhar123/ciba-server/grant"
	"log"
	"net/http"
)

type AuthorizationServerProvider interface {
	AddGrant(grant grant.GrantTypeInterface)
	HandleAuthentication(r http.Request)
}

type AuthorizationServer struct {
	grants map[string]grant.GrantTypeInterface
}

func (as *AuthorizationServer) AddGrant(gt grant.GrantTypeInterface) {
	_, exist := as.grants[gt.GetIdentifier()]
	if !exist {
		as.grants[gt.GetIdentifier()] = gt
		log.Printf("added grant type: %s\n", gt.GetIdentifier())
	}
}