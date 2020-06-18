package repository

import "github.com/adisazhar123/ciba-server/domain"

// to store access token
type AccessTokenRepositoryInterface interface {

}

// to store ciba session
type CibaSessionRepositoryInterface interface {
	create(cibaSession domain.CibaSession)
}

// to store client application
type ClientApplicationRepositoryInterface interface {
	register(clientApp domain.ClientApplicationInterface)
}

// to store public & private key
type KeyRepositoryInterface interface {

}