package repository

import "github.com/jmoiron/sqlx"

type AuthRepo struct {
}

func NewAuthRepo(db *sqlx.DB) *AuthRepo {
	return &AuthRepo{}
}
