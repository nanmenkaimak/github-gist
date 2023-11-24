package repository

import (
	"github.com/google/uuid"
	"github.com/nanmenkaimak/github-gist/internal/admin/database/dbpostgres"
	"github.com/nanmenkaimak/github-gist/internal/admin/entity"
)

type Gist interface {
	GetOtherAllGists(sort string, direction string) (*[]entity.GistRequest, error)
	GetGistByID(gistID uuid.UUID) (*entity.GistRequest, error)
	DeleteGistByID(id uuid.UUID) error
}

type User interface {
	UpdateUser(updatedUser entity.User) error
	GetUserByUsername(username string) (*entity.User, error)
	GetAllUsers() (*[]entity.User, error)
}

type GistRepo struct {
	main    *dbpostgres.Db
	replica *dbpostgres.Db
}

type UserRepo struct {
	main    *dbpostgres.Db
	replica *dbpostgres.Db
}

func NewGistRepository(connMain *dbpostgres.Db, connReplica *dbpostgres.Db) *GistRepo {
	return &GistRepo{
		main:    connMain,
		replica: connReplica,
	}
}

func NewUserRepository(connMain *dbpostgres.Db, connReplica *dbpostgres.Db) *UserRepo {
	return &UserRepo{
		main:    connMain,
		replica: connReplica,
	}
}
