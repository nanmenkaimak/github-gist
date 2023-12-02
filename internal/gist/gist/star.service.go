package gist

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/nanmenkaimak/github-gist/internal/gist/entity"
)

func (a *Service) StarGist(ctx context.Context, request entity.Star) error {
	err := a.repo.StarGist(request)
	if err != nil {
		return fmt.Errorf("StarGist request err: %v", err)
	}
	return nil
}

func (a *Service) GetStaredGists(ctx context.Context, request OtherGistRequest) (*[]entity.GistRequest, error) {
	user, err := a.userGrpcTransport.GetUserByUsername(ctx, request.Username)
	if err != nil {
		return nil, fmt.Errorf("GetUserByUsername request err: %v", err)
	}

	userID, err := uuid.Parse(user.Id)
	if err != nil {
		return nil, fmt.Errorf("parse uuid err: %v", err)
	}

	gists, err := a.repo.GetStarredGists(userID)
	if err != nil {
		return nil, fmt.Errorf("getting starred gists err: %v", err)
	}

	return &gists, nil
}

func (a *Service) DeleteStar(ctx context.Context, request DeleteRequest) error {
	err := a.repo.DeleteStar(request.GistID, request.UserID)
	if err != nil {
		return fmt.Errorf("delete star err: %v", err)
	}
	return nil
}
