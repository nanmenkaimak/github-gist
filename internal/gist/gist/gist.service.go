package gist

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/nanmenkaimak/github-gist/internal/gist/entity"
)

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

func (a *Service) GetAllGists(ctx context.Context, request GetAllGistsRequest) (*[]entity.GistRequest, error) {
	if request.Sort == "" {
		request.Sort = "created_at"
	}
	if request.Direction == "" {
		request.Direction = "desc"
	}
	allGists, err := a.repo.GetOtherAllGists(request.Sort, request.Direction)
	if err != nil {
		return nil, fmt.Errorf("getting all gists err: %v", err)
	}

	return &allGists, nil
}

func (a *Service) GetGistByID(ctx context.Context, request GetGistRequest) (*entity.GistRequest, error) {
	user, err := a.userGrpcTransport.GetUserByUsername(ctx, request.Username)
	if err != nil {
		return nil, fmt.Errorf("GetUserByUsername request err: %v", err)
	}

	userID, err := uuid.Parse(user.Id)
	if err != nil {
		return nil, fmt.Errorf("parse uuid err: %v", err)
	}

	ownGist := false
	if userID == request.UserID {
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
	user, err := a.userGrpcTransport.GetUserByUsername(ctx, request.Username)
	if err != nil {
		return nil, fmt.Errorf("GetUserByUsername request err: %v", err)
	}

	userID, err := uuid.Parse(user.Id)
	if err != nil {
		return nil, fmt.Errorf("parse uuid err: %v", err)
	}

	ownGist := false
	if userID == request.UserID {
		ownGist = true
	}

	gists, err := a.repo.GetAllGistsOfUser(userID, ownGist, request.Searching)
	if err != nil {
		return nil, fmt.Errorf("getting all gists of user err: %v", err)
	}

	return &gists, err
}

func (a *Service) GetGistsByVisibility(ctx context.Context, request GetGistRequest) (*[]entity.GistRequest, error) {
	user, err := a.userGrpcTransport.GetUserByUsername(ctx, request.Username)
	if err != nil {
		return nil, fmt.Errorf("GetUserByUsername request err: %v", err)
	}

	userID, err := uuid.Parse(user.Id)
	if err != nil {
		return nil, fmt.Errorf("parse uuid err: %v", err)
	}

	if userID != request.UserID && request.Visibility == false {
		return nil, fmt.Errorf("different user err: %v", err)
	}

	gists, err := a.repo.GetGistsByVisibility(request.UserID, request.Visibility)
	if err != nil {
		return nil, fmt.Errorf("getting gists by visibility err: %v", err)
	}

	return &gists, nil
}

func (a *Service) UpdateGistByID(ctx context.Context, request UpdateGistRequest) error {
	user, err := a.userGrpcTransport.GetUserByUsername(ctx, request.Username)
	if err != nil {
		return fmt.Errorf("GetUserByUsername request err: %v", err)
	}

	userID, err := uuid.Parse(user.Id)
	if err != nil {
		return fmt.Errorf("parse uuid err: %v", err)
	}

	if userID != request.UserID {
		return fmt.Errorf("it is not your gist err: %v", err)
	}

	err = a.repo.UpdateGistByID(request.Gist)
	if err != nil {
		return fmt.Errorf("update gist err: %v", err)
	}
	return nil
}

func (a *Service) DeleteGistByID(ctx context.Context, request GetGistRequest) error {
	user, err := a.userGrpcTransport.GetUserByUsername(ctx, request.Username)
	if err != nil {
		return fmt.Errorf("GetUserByUsername request err: %v", err)
	}

	userID, err := uuid.Parse(user.Id)
	if err != nil {
		return fmt.Errorf("parse uuid err: %v", err)
	}

	if userID != request.UserID {
		return fmt.Errorf("it is not your gist err: %v", err)
	}

	err = a.repo.DeleteGistByID(request.GistID)
	if err != nil {
		return fmt.Errorf("delete gist err: %v", err)
	}
	return nil
}
