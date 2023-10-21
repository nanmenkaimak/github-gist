package repository

import (
	"github.com/nanmenkaimak/github-gist/internal/auth/database/dbpostgres"
	"github.com/nanmenkaimak/github-gist/internal/auth/entitiy"
)

type Repository interface {
	UserTokenRepository
}

type UserTokenRepository interface {
	CreateUserToken(userToken entitiy.UserToken) error
	UpdateUserToken(userToken entitiy.UserToken) error
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
