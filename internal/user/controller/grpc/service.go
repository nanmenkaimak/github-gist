package grpc

import (
	"context"
	"fmt"
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
