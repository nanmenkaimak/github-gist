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
	repo              repository.Repository
	userGrpcTransport *transport.UserGrpcTransport
}

func NewGistService(repo repository.Repository, userGrpcTransport *transport.UserGrpcTransport) UseCase {
	return &Service{
		repo:              repo,
		userGrpcTransport: userGrpcTransport,
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

	if userID != request.UserID {
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

	err = a.repo.DeleteStar(request.GistID, request.UserID)
	if err != nil {
		return fmt.Errorf("delete star err: %v", err)
	}
	return nil
}

func (a *Service) ForkGist(ctx context.Context, request ForkRequest) (*ForkGistResponse, error) {
	user, err := a.userGrpcTransport.GetUserByUsername(ctx, request.Username)
	if err != nil {
		return nil, fmt.Errorf("GetUserByUsername request err: %v", err)
	}

	userID, err := uuid.Parse(user.Id)
	if err != nil {
		return nil, fmt.Errorf("parse uuid err: %v", err)
	}

	if userID != request.UserID {
		return nil, fmt.Errorf("it is not your gist err: %v", err)
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

func (a *Service) CreateComment(ctx context.Context, newComment entity.Comment) error {
	user, err := a.userGrpcTransport.GetUserByID(ctx, newComment.UserID.String())
	if err != nil {
		return fmt.Errorf("GetUserByID request err: %v", err)
	}
	newComment.Username = user.Username
	err = a.repo.CreateComment(newComment)
	if err != nil {
		return err
	}
	return nil
}

func (a *Service) GetCommentsOfGist(ctx context.Context, request GetGistRequest) (*[]entity.Comment, error) {
	gists, err := a.repo.GetAllCommentsOfGist(request.GistID)
	if err != nil {
		return nil, err
	}
	return &gists, err
}

func (a *Service) DeleteComment(ctx context.Context, request DeleteRequest) error {
	user, err := a.userGrpcTransport.GetUserByUsername(ctx, request.Username)
	if err != nil {
		return fmt.Errorf("GetUserByUsername request err: %v", err)
	}

	userID, err := uuid.Parse(user.Id)
	if err != nil {
		return fmt.Errorf("parse uuid err: %v", err)
	}

	if userID != request.UserID {
		return fmt.Errorf("it is not your comment err: %v", err)
	}

	err = a.repo.DeleteComment(request.CommentID)
	if err != nil {
		return fmt.Errorf("delete comment err: %v", err)
	}
	return nil
}

func (a *Service) UpdateComment(ctx context.Context, request UpdateCommentRequest) error {
	user, err := a.userGrpcTransport.GetUserByUsername(ctx, request.Username)
	if err != nil {
		return fmt.Errorf("GetUserByUsername request err: %v", err)
	}

	userID, err := uuid.Parse(user.Id)
	if err != nil {
		return fmt.Errorf("parse uuid err: %v", err)
	}

	if userID != request.UserID {
		return fmt.Errorf("it is not your comment err: %v", err)
	}

	updatedComment := entity.Comment{
		ID:   request.CommentID,
		Text: request.Text,
	}

	err = a.repo.UpdateComment(updatedComment)
	if err != nil {
		return fmt.Errorf("update comment err: %v", err)
	}
	return nil
}
