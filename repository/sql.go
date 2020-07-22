package repository

import (
	"github.com/adisazhar123/go-ciba/domain"
	"github.com/jmoiron/sqlx"
)

type ClientApplicationSQLRepository struct {
	db sqlx.DB
}

func (repo *ClientApplicationSQLRepository) Register(clientApp domain.ClientApplicationInterface) {

}
