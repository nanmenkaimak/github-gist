package grpc

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/nanmenkaimak/github-gist/internal/user/entity"
	"github.com/nanmenkaimak/github-gist/internal/user/repository"
	pb "github.com/nanmenkaimak/github-gist/pkg/protobuf/userservice/gw"
	"go.uber.org/zap"
)

type Service struct {
	pb.UnimplementedUserServiceServer
	logger *zap.SugaredLogger
	repo   repository.Repository
}

func NewService(logger *zap.SugaredLogger, repo repository.Repository) *Service {
	return &Service{
		logger: logger,
		repo:   repo,
	}
}

func (s *Service) CreateUser(ctx context.Context, request *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	userRequest := entity.User{
		FirstName:   request.User.FirstName,
		LastName:    request.User.LastName,
		Username:    request.User.Username,
		Email:       request.User.Email,
		Password:    request.User.Password,
		IsConfirmed: request.User.IsConfirmed,
	}
	userID, err := s.repo.CreateUser(userRequest)
	if err != nil {
		s.logger.Errorf("failed to CreateUser err: %v", err)
		return nil, fmt.Errorf("CreateUser err: %w", err)
	}

	return &pb.CreateUserResponse{
		Id: userID.String(),
	}, nil
}

func (s *Service) GetUserByUsername(ctx context.Context, request *pb.GetUserByUsernameRequest) (*pb.GetUserByUsernameResponse, error) {
	user, err := s.repo.GetUserByUsername(request.Username)
	if err != nil {
		s.logger.Errorf("failed to GetUserByLogin err: %v", err)
		return nil, fmt.Errorf("GetUserByLogin err: %w", err)
	}

	return &pb.GetUserByUsernameResponse{
		Result: &pb.User{
			Id:          user.ID.String(),
			FirstName:   user.FirstName,
			LastName:    user.LastName,
			Username:    user.Username,
			Email:       user.Email,
			Password:    user.Password,
			IsConfirmed: user.IsConfirmed,
		},
	}, nil
}

func (s *Service) ConfirmUser(ctx context.Context, request *pb.ConfirmUserRequest) (*pb.ConfirmUserResponse, error) {
	err := s.repo.ConfirmUser(request.GetEmail())
	if err != nil {
		s.logger.Errorf("failed to ConfirmUser err: %v", err)
		return nil, fmt.Errorf("ConfirmUser err: %w", err)
	}

	return &pb.ConfirmUserResponse{}, nil
}

func (s *Service) GetUserByID(ctx context.Context, request *pb.GetUserByIDRequest) (*pb.GetUserByIDResponse, error) {
	user, err := s.repo.GetUserByID(request.GetId())
	if err != nil {
		s.logger.Errorf("failed to GetUserByID err: %v", err)
		return nil, fmt.Errorf("GetUserByID err: %v", err)
	}
	return &pb.GetUserByIDResponse{
		User: &pb.User{
			Id:          user.ID.String(),
			FirstName:   user.FirstName,
			LastName:    user.LastName,
			Username:    user.Username,
			Email:       user.Email,
			Password:    user.Password,
			IsConfirmed: user.IsConfirmed,
		},
	}, nil
}

func (s *Service) FollowUser(ctx context.Context, request *pb.FollowUserRequest) (*pb.FollowUserResponse, error) {
	followingID, err := uuid.Parse(request.FollowingId)
	if err != nil {
		return nil, fmt.Errorf("converting string to uuid err: %v", err)
	}
	followerID, err := uuid.Parse(request.FollowerId)
	if err != nil {
		return nil, fmt.Errorf("converting string to uuid err: %v", err)
	}
	newFollow := entity.Follower{
		FollowerID:  followerID,
		FollowingID: followingID,
	}
	err = s.repo.FollowUser(newFollow)
	if err != nil {
		s.logger.Errorf("failed to FollowUser err: %v", err)
		return nil, fmt.Errorf("FollowUser err: %v", err)
	}

	return &pb.FollowUserResponse{}, nil
}

func (s *Service) UnfollowUser(ctx context.Context, request *pb.UnfollowUserRequest) (*pb.UnfollowUserResponse, error) {
	followingID, err := uuid.Parse(request.FollowingId)
	if err != nil {
		return nil, fmt.Errorf("converting string to uuid err: %v", err)
	}
	followerID, err := uuid.Parse(request.FollowerId)
	if err != nil {
		return nil, fmt.Errorf("converting string to uuid err: %v", err)
	}
	newUnfollow := entity.Follower{
		FollowerID:  followerID,
		FollowingID: followingID,
	}
	err = s.repo.UnfollowUser(newUnfollow)
	if err != nil {
		s.logger.Errorf("failed to UnfollowUser err: %v", err)
		return nil, fmt.Errorf("UnfollowUser err: %v", err)
	}

	return &pb.UnfollowUserResponse{}, nil
}

func (s *Service) GetAllFollowers(ctx context.Context, request *pb.GetAllFollowersRequest) (*pb.GetAllFollowersResponse, error) {
	users, err := s.repo.GetAllFollowers(request.UserId)
	if err != nil {
		s.logger.Errorf("failed to GetAllFollowers err: %v", err)
		return nil, fmt.Errorf("GetAllFollowers err: %v", err)
	}

	var responseUsers []*pb.User

	for i := 0; i < len(users); i++ {
		follower := pb.User{
			Id:          users[i].ID.String(),
			FirstName:   users[i].FirstName,
			LastName:    users[i].LastName,
			Username:    users[i].Username,
			Email:       users[i].Email,
			Password:    users[i].Password,
			IsConfirmed: users[i].IsConfirmed,
		}
		responseUsers = append(responseUsers, &follower)
	}

	return &pb.GetAllFollowersResponse{
		Followers: responseUsers,
	}, nil
}

func (s *Service) GetAllFollowings(ctx context.Context, request *pb.GetAllFollowingsRequest) (*pb.GetAllFollowingsResponse, error) {
	users, err := s.repo.GetAllFollowings(request.UserId)
	if err != nil {
		s.logger.Errorf("failed to GetAllFollowings err: %v", err)
		return nil, fmt.Errorf("GetAllFollowings err: %v", err)
	}

	var responseUsers []*pb.User

	for i := 0; i < len(users); i++ {
		following := pb.User{
			Id:          users[i].ID.String(),
			FirstName:   users[i].FirstName,
			LastName:    users[i].LastName,
			Username:    users[i].Username,
			Email:       users[i].Email,
			Password:    users[i].Password,
			IsConfirmed: users[i].IsConfirmed,
		}
		responseUsers = append(responseUsers, &following)
	}

	return &pb.GetAllFollowingsResponse{
		Followings: responseUsers,
	}, nil
}
