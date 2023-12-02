package gist

import (
	"context"
	"fmt"

	"github.com/nanmenkaimak/github-gist/internal/gist/entity"
)

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
	err := a.repo.DeleteComment(request.CommentID, request.UserID)
	if err != nil {
		return fmt.Errorf("delete comment err: %v", err)
	}
	return nil
}

func (a *Service) UpdateComment(ctx context.Context, request UpdateCommentRequest) error {
	updatedComment := entity.Comment{
		ID:     request.CommentID,
		UserID: request.UserID,
		Text:   request.Text,
	}

	err := a.repo.UpdateComment(updatedComment)
	if err != nil {
		return fmt.Errorf("update comment err: %v", err)
	}
	return nil
}
