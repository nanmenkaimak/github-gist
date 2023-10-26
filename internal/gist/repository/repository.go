package repository

import (
	"github.com/google/uuid"
	"github.com/nanmenkaimak/github-gist/internal/gist/database/dbpostgres"
	"github.com/nanmenkaimak/github-gist/internal/gist/entity"
)

type Repository interface {
	GistRepository
	ForkRepository
	Comment
	Star
}

type GistRepository interface {
	CreateGist(request entity.GistRequest) (uuid.UUID, error)
	GetOtherAllGists(sort string, direction string) ([]entity.GistRequest, error)
	GetGistByID(gistID uuid.UUID, ownGist bool) (entity.GistRequest, error)
	GetAllGistsOfUser(userID uuid.UUID, ownGists bool, searchingStr string) ([]entity.GistRequest, error)
	GetGistsByVisibility(userID uuid.UUID, visibility bool) ([]entity.GistRequest, error)
	UpdateGistByID(updatedGist entity.GistRequest) error
	DeleteGistByID(id uuid.UUID) error
}

type ForkRepository interface {
}

type Comment interface {
}

type Star interface {
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
