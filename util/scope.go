package util

import (
	"github.com/adisazhar123/ciba-server/domain"
	"strings"
)

type ScopeUtil struct {

}

func find(slice []string, val string) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}

func (ss *ScopeUtil) ScopeExist(clientApp *domain.ClientApplication, scope string) bool {
	registeredScope := strings.Split(clientApp.Scope, " ")
	requestedScope := strings.Split(scope, " ")

	if len(registeredScope) < len(requestedScope) {
		return false
	}

	for _, v := range requestedScope {
		if find(registeredScope, v) {
			return true
		}
	}
	return false
}