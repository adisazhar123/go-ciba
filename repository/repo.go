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
	Register(clientApp domain.ClientApplication) error
	FindById(id string) *domain.ClientApplication
}

// to store public & private key
type KeyRepositoryInterface interface {

}

type UserAccountRepositoryInterface interface {
	FindById(id string) (*domain.UserAccount, error)
}