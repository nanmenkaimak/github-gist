package gist

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/nanmenkaimak/github-gist/internal/gist/entity"
	"github.com/nanmenkaimak/github-gist/internal/gist/repository"
	"github.com/nanmenkaimak/github-gist/internal/gist/transport"
)

type Service struct {
	repo          repository.Repository
	userTransport *transport.UserTransport
}

func NewGistService(repo repository.Repository, userTransport *transport.UserTransport) UseCase {
	return &Service{
		repo:          repo,
		userTransport: userTransport,
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

func (a *Service) GetGistByID(ctx context.Context, request GetGistRequest) (*entity.GistRequest, error) {
	user, err := a.userTransport.GetUser(ctx, request.Username)
	if err != nil {
		return nil, fmt.Errorf("GetUser request err: %v", err)
	}
	ownGist := false
	if user.ID == request.UserID {
		ownGist = true
	}
	gist, err := a.repo.GetGistByID(request.GistID, ownGist)
	if err != nil {
		return nil, fmt.Errorf("getting gist err: %v", err)
	}

	if gist.Gist.ID == uuid.Nil {
		return nil, nil
	}

	return &gist, nil
}

func (a *Service) GetAllGistsOfUser(ctx context.Context, request GetGistRequest) (*[]entity.GistRequest, error) {
	user, err := a.userTransport.GetUser(ctx, request.Username)
	if err != nil {
		return nil, fmt.Errorf("GetUser request err: %v", err)
	}
	ownGist := false
	if user.ID == request.UserID {
		ownGist = true
	}
	gists, err := a.repo.GetAllGistsOfUser(user.ID, ownGist)
	if err != nil {
		return nil, fmt.Errorf("getting all gists of user err: %v", err)
	}

	return &gists, err
}
