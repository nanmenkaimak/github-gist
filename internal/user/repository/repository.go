package repository

import "github.com/nanmenkaimak/github-gist/internal/user/database/dbpostgres"

type Repository interface {
	UserRepository
}

type UserRepository interface {
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
