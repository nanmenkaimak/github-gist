package gist

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/nanmenkaimak/github-gist/internal/gist/entity"
	"github.com/nanmenkaimak/github-gist/internal/gist/repository"
)

type Service struct {
	repo repository.Repository
}

func NewGistService(repo repository.Repository) UseCase {
	return &Service{
		repo: repo,
	}
}

func (a *Service) CreateGist(ctx context.Context, gistRequest entity.GistRequest) (*CreateGistResponse, error) {
	gistID, err := a.repo.CreateGist(gistRequest)
	if err != nil {
		return nil, fmt.Errorf("creating gist err: %v", err)
	}

	response := &CreateGistResponse{
		GistID: gistID,
	}

	return response, nil
}

func (a *Service) GetAllGists(ctx context.Context) (*[]entity.GistRequest, error) {
	allGists, err := a.repo.GetOtherAllGists()
	if err != nil {
		return nil, fmt.Errorf("getting all gists err: %v", err)
	}

	return &allGists, nil
}

func (a *Service) GetGistByID(ctx context.Context, request GetGistByIDRequest) (*entity.GistRequest, error) {
	gist, err := a.repo.GetOtherGistByID(request.GistID)
	if err != nil {
		return nil, fmt.Errorf("getting gist err: %v", err)
	}

	if gist.Gist.ID == uuid.Nil {
		return nil, nil
	}

	return &gist, nil
}
