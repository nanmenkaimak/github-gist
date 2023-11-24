package admin

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/nanmenkaimak/github-gist/internal/admin/entity"
)

func (a *Service) GetAllGists(ctx context.Context, request GetAllGistsRequest) (*[]entity.GistRequest, error) {
	if request.Sort == "" {
		request.Sort = "created_at"
	}
	if request.Direction == "" {
		request.Direction = "desc"
	}
	allGists, err := a.gistRepo.GetOtherAllGists(request.Sort, request.Direction)
	if err != nil {
		return nil, fmt.Errorf("getting all gists err: %v", err)
	}

	return allGists, nil
}

func (a *Service) GetGistByID(ctx context.Context, request GetGistRequest) (*entity.GistRequest, error) {
	gist, err := a.gistRepo.GetGistByID(request.GistID)
	if err != nil {
		return nil, fmt.Errorf("getting gist err: %v", err)
	}

	if gist.Gist.ID == uuid.Nil {
		return nil, nil
	}

	return gist, nil
}

func (a *Service) DeleteGistByID(ctx context.Context, request GetGistRequest) error {
	err := a.gistRepo.DeleteGistByID(request.GistID)
	if err != nil {
		return fmt.Errorf("delete gist err: %v", err)
	}
	return nil
}
