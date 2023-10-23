package gist

import (
	"context"
	"github.com/nanmenkaimak/github-gist/internal/gist/entity"
)

type UseCase interface {
	CreateGist(ctx context.Context, gistRequest entity.GistRequest) (*CreateGistResponse, error)
	GetAllGists(ctx context.Context) (*[]entity.GistRequest, error)
	GetGistByID(ctx context.Context, request GetGistRequest) (*entity.GistRequest, error)
	GetAllGistsOfUser(ctx context.Context, request GetGistRequest) (*[]entity.GistRequest, error)
	UpdateGistByID(ctx context.Context, request UpdateGistRequest) error
}
