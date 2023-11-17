package gist

import (
	"context"
	"fmt"
	"github.com/google/uuid"
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
