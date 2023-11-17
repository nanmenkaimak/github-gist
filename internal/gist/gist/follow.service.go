package gist

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/nanmenkaimak/github-gist/internal/gist/entity"
)

func (a *Service) FollowUser(ctx context.Context, request FollowRequest) error {
	user, err := a.userGrpcTransport.GetUserByUsername(ctx, request.Username)
	if err != nil {
		return fmt.Errorf("GetUserByUsername request err: %v", err)
	}

	followingID, err := uuid.Parse(user.Id)
	if err != nil {
		return fmt.Errorf("converting string to uuid err: %v", err)
	}

	request.FollowingID = followingID

	_, err = a.userGrpcTransport.FollowUser(ctx, request.FollowerID.String(), request.FollowingID.String())
	if err != nil {
		return fmt.Errorf("FollowUser request err: %v", err)
	}
	return nil
}

func (a *Service) UnfollowUser(ctx context.Context, request FollowRequest) error {
	user, err := a.userGrpcTransport.GetUserByUsername(ctx, request.Username)
	if err != nil {
		return fmt.Errorf("GetUserByUsername request err: %v", err)
	}

	followingID, err := uuid.Parse(user.Id)
	if err != nil {
		return fmt.Errorf("converting string to uuid err: %v", err)
	}

	request.FollowingID = followingID

	_, err = a.userGrpcTransport.UnfollowUser(ctx, request.FollowerID.String(), request.FollowingID.String())
	if err != nil {
		return fmt.Errorf("FollowUser request err: %v", err)
	}
	return nil
}

func (a *Service) GetAllFollowers(ctx context.Context, username string) (*[]entity.UserResponse, error) {
	user, err := a.userGrpcTransport.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("GetUserByUsername request err: %v", err)
	}

	resp, err := a.userGrpcTransport.GetAllFollowers(ctx, user.Id)
	if err != nil {
		return nil, fmt.Errorf("GetAllFollowers request err: %v", err)
	}

	return resp, nil
}

func (a *Service) GetAllFollowings(ctx context.Context, username string) (*[]entity.UserResponse, error) {
	user, err := a.userGrpcTransport.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("GetUserByUsername request err: %v", err)
	}

	resp, err := a.userGrpcTransport.GetAllFollowings(ctx, user.Id)
	if err != nil {
		return nil, fmt.Errorf("GetAllFollowings request err: %v", err)
	}

	return resp, nil
}

func (a *Service) GetUserInfo(ctx context.Context, username string) (entity.UserResponse, error) {
	user, err := a.userGrpcTransport.GetUserByUsername(ctx, username)
	if err != nil {
		return entity.UserResponse{}, fmt.Errorf("GetUserByUsername request err: %v", err)
	}

	responseUser := entity.UserResponse{
		ID:          user.Id,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		Username:    user.Username,
		Email:       user.Email,
		Password:    user.Password,
		IsConfirmed: user.IsConfirmed,
	}

	return responseUser, nil
}
