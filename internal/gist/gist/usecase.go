package gist

import (
	"context"
	"github.com/nanmenkaimak/github-gist/internal/gist/entity"
)

type UseCase interface {
	Gist
	Star
	Fork
	Comment
}

type Gist interface {
	CreateGist(ctx context.Context, gistRequest entity.GistRequest) (*CreateGistResponse, error)
	GetAllGists(ctx context.Context, request GetAllGistsRequest) (*[]entity.GistRequest, error)
	GetGistByID(ctx context.Context, request GetGistRequest) (*entity.GistRequest, error)
	GetAllGistsOfUser(ctx context.Context, request GetGistRequest) (*[]entity.GistRequest, error)
	GetGistsByVisibility(ctx context.Context, request GetGistRequest) (*[]entity.GistRequest, error)
	UpdateGistByID(ctx context.Context, request UpdateGistRequest) error
	DeleteGistByID(ctx context.Context, request GetGistRequest) error
}

type Star interface {
	StarGist(ctx context.Context, request entity.Star) error
	GetStaredGists(ctx context.Context, request OtherGistRequest) (*[]entity.GistRequest, error)
	DeleteStar(ctx context.Context, request DeleteRequest) error //
}

type Fork interface {
	ForkGist(ctx context.Context, request ForkRequest) (*ForkGistResponse, error)
	GetForkedGists(ctx context.Context, request GetGistRequest) (*[]entity.GistRequest, error)
}

type Comment interface {
	CreateComment(ctx context.Context, newComment entity.Comment) error
	GetCommentsOfGist(ctx context.Context, request GetGistRequest) (*[]entity.Comment, error)
	DeleteComment()
	UpdateComment()
}
