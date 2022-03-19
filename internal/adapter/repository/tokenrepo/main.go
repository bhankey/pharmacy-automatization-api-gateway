package tokenrepo

import "github.com/jmoiron/sqlx"

type TokenRepo struct {
	master *sqlx.DB
	slave  *sqlx.DB
}

func NewTokenRepo(master *sqlx.DB, slave *sqlx.DB) *TokenRepo {
	return &TokenRepo{
		master: master,
		slave:  slave,
	}
}
