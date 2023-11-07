package repository

import (
	"github.com/google/uuid"
	"github.com/nanmenkaimak/github-gist/internal/user/database/dbpostgres"
	"github.com/nanmenkaimak/github-gist/internal/user/entity"
)

type Repository interface {
	UserRepository
}

type UserRepository interface {
	CreateUser(newUser entity.User) (uuid.UUID, error)
	GetUserByUsername(username string) (*entity.User, error)
	ConfirmUser(email string) error
	GetUserByID(userID string) (*entity.User, error)
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
