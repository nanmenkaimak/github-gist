package repository

import (
	"github.com/google/uuid"
	"github.com/nanmenkaimak/github-gist/internal/gist/database/dbpostgres"
	"github.com/nanmenkaimak/github-gist/internal/gist/entity"
)

type Repository interface {
	GistRepository
	ForkRepository
	CommentRepository
	StarRepository
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
	ForkGist(newFork entity.Fork) error
	GetForkedGistByUser(userID uuid.UUID, ownGists bool) ([]entity.GistRequest, error)
	DeleteFork(id uuid.UUID) error
}

type CommentRepository interface {
	CreateComment(newComment entity.Comment) error
	GetAllCommentsOfGist(gistID uuid.UUID) ([]entity.Comment, error)
	DeleteComment(id uuid.UUID) error
	UpdateComment(updatedComment entity.Comment) error
}

type StarRepository interface {
	StarGist(newStar entity.Star) error
	GetStarredGists(userID uuid.UUID) ([]entity.GistRequest, error)
	GetAllStargazers()
	DeleteStar(gistID uuid.UUID, userID uuid.UUID) error
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
