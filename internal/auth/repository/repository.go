package repository

import (
	"github.com/nanmenkaimak/github-gist/internal/auth/database/dbpostgres"
	"github.com/nanmenkaimak/github-gist/internal/auth/entity"
)

type Repository interface {
	UserTokenRepository
}

type UserTokenRepository interface {
	CreateUserToken(userToken entity.UserToken) error
	UpdateUserToken(userToken entity.UserToken) error
}

type Repo struct {
	main    *dbpostgres.Db
	replica *dbpostgres.Db
}

func NewRepository(connMain *dbpostgres.Db, connReplica *dbpostgres.Db) *Repo {
	return &Repo{
		main:    connMain,
		replica: connReplica,
	}
}
