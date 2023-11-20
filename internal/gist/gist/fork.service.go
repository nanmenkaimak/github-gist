package gist

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/nanmenkaimak/github-gist/internal/gist/entity"
)

func (a *Service) ForkGist(ctx context.Context, request ForkRequest) (*ForkGistResponse, error) {
	user, err := a.userGrpcTransport.GetUserByUsername(ctx, request.Username)
	if err != nil {
		return nil, fmt.Errorf("GetUserByUsername request err: %v", err)
	}

	userID, err := uuid.Parse(user.Id)
	if err != nil {
		return nil, fmt.Errorf("parse uuid err: %v", err)
	}

	if userID == request.UserID {
		return nil, fmt.Errorf("it is your gist err: %v", err)
	}

	gist, err := a.repo.GetGistByID(request.GistID, false)
	if err != nil {
		return nil, fmt.Errorf("getting gist err: %v", err)
	}

	var files []entity.File

	for i := 0; i < len(gist.Files); i++ {
		var file entity.File
		file.Name = gist.Files[i].Name
		file.Code = gist.Files[i].Code
		files = append(files, file)
	}
	newForkedGist := entity.GistRequest{
		Gist: entity.Gist{
			UserID:      request.UserID,
			Name:        gist.Gist.UserID.String(),
			Description: gist.Gist.Description,
			Visible:     gist.Gist.Visible,
			IsForked:    true,
		},
		Commit: entity.Commit{
			Comment: gist.Commit.Comment,
		},
		Files: files,
	}

	forkedGistID, err := a.repo.CreateGist(newForkedGist)
	if err != nil {
		return nil, fmt.Errorf("creating gist err: %v", err)
	}

	forkRequest := entity.Fork{
		GistID:    gist.Gist.ID,
		NewGistID: forkedGistID,
	}

	err = a.repo.ForkGist(forkRequest)
	if err != nil {
		return nil, fmt.Errorf("creating fork err: %v", err)
	}

	return &ForkGistResponse{
		GistID: forkedGistID,
	}, nil
}

func (a *Service) GetForkedGists(ctx context.Context, request GetGistRequest) (*[]entity.GistRequest, error) {
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

	gists, err := a.repo.GetForkedGistByUser(userID, ownGist)
	if err != nil {
		return nil, fmt.Errorf("getting forked gists err: %v", err)
	}

	return &gists, nil
}
