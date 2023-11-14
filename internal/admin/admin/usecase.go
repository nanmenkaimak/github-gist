package admin

import (
	"context"

	"github.com/nanmenkaimak/github-gist/internal/admin/entity"
)

type UseCase interface {
	Gist
	User
}

type Gist interface {
	GetAllGists(ctx context.Context, request GetAllGistsRequest) (*[]entity.GistRequest, error)
	GetGistByID(ctx context.Context, request GetGistRequest) (*entity.GistRequest, error)
	DeleteGistByID(ctx context.Context, request GetGistRequest) error
}

type User interface {
	UpdateUserByUsername(ctx context.Context, request UpdateUserRequest) error
	GetAllUsers(ctx context.Context, request GetUserRequest) (*[]entity.User, error)
	GetUserByUsername(ctx context.Context, request GetUserRequest) (*entity.User, error)
}
