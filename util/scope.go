package util

import (
	"strings"
)

type ScopeUtil struct {
}

func (ss *ScopeUtil) ScopeExist(clientAppScope string, scope string) bool {
	registeredScope := strings.Split(clientAppScope, " ")
	requestedScope := strings.Split(scope, " ")

	if len(registeredScope) < len(requestedScope) {
		return false
	}

	mapRegistered := make(map[string]struct{})

	for _, v := range registeredScope {
		mapRegistered[v] = struct{}{}
	}

	for _, v := range requestedScope {
		if _, exist := mapRegistered[v]; !exist {
			return false
		}
	}

	return true
}
