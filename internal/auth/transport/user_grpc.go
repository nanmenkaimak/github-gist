package transport

import (
	"context"
	"fmt"
	"github.com/nanmenkaimak/github-gist/internal/auth/config"
	"github.com/nanmenkaimak/github-gist/internal/auth/entitiy"
	pb "github.com/nanmenkaimak/github-gist/pkg/protobuf/userservice/gw"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type UserGrpcTransport struct {
	config config.UserGrpcTransport
	client pb.UserServiceClient
}

func NewUserGrpcTransport(config config.UserGrpcTransport) *UserGrpcTransport {
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	conn, _ := grpc.Dial(config.Host, opts...)

	client := pb.NewUserServiceClient(conn)

	return &UserGrpcTransport{
		config: config,
		client: client,
	}
}

func (t *UserGrpcTransport) GetUserByUsername(ctx context.Context, username string) (*pb.User, error) {
	resp, err := t.client.GetUserByUsername(ctx, &pb.GetUserByUsernameRequest{
		Username: username,
	})
	if err != nil {
		return nil, fmt.Errorf("cannot GetUserByUsername: %w", err)
	}

	if resp == nil {
		return nil, fmt.Errorf("not found")
	}

	return resp.Result, nil
}

func (t *UserGrpcTransport) CreateUser(ctx context.Context, newUser entitiy.RegisterUserRequest) (*pb.CreateUserResponse, error) {
	resp, err := t.client.CreateUser(ctx, &pb.CreateUserRequest{
		User: &pb.User{
			FirstName:   newUser.FirstName,
			LastName:    newUser.LastName,
			Username:    newUser.Username,
			Email:       newUser.Email,
			Password:    newUser.Password,
			IsConfirmed: newUser.IsConfirmed,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("cannot CreateUser: %w", err)
	}
	return resp, nil
}

func (t *UserGrpcTransport) ConfirmUser(ctx context.Context, email string) (*pb.ConfirmUserResponse, error) {
	resp, err := t.client.ConfirmUser(ctx, &pb.ConfirmUserRequest{
		Email: email,
	})

	if err != nil {
		return nil, fmt.Errorf("cannot ConfirmUser: %w", err)
	}
	return resp, nil
}

func (t *UserGrpcTransport) UpdateUser(ctx context.Context, updatedUser entitiy.RegisterUserRequest) (*pb.UpdateUserResponse, error) {
	resp, err := t.client.UpdateUser(ctx, &pb.UpdateUserRequest{
		User: &pb.User{
			FirstName: updatedUser.FirstName,
			LastName:  updatedUser.LastName,
			Username:  updatedUser.Username,
			Email:     updatedUser.Email,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("cannot UpdateUser: %w", err)
	}
	return resp, nil
}
