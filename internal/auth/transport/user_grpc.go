package transport

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"io"

	"github.com/nanmenkaimak/github-gist/internal/auth/config"
	"github.com/nanmenkaimak/github-gist/internal/auth/entity"
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

func (t *UserGrpcTransport) CreateUser(ctx context.Context, newUser entity.RegisterUserRequest) (*pb.CreateUserResponse, error) {
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

func (t *UserGrpcTransport) UpdateUser(ctx context.Context, updatedUser entity.RegisterUserRequest) (*pb.UpdateUserResponse, error) {
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

func (t *UserGrpcTransport) UpdatePassword(ctx context.Context, email string, newPassword string) (*pb.UpdatePasswordResponse, error) {
	resp, err := t.client.UpdatePassword(ctx, &pb.UpdatePasswordRequest{
		NewPassword: newPassword,
		Email:       email,
	})
	if err != nil {
		return nil, fmt.Errorf("cannot UpdatePassword: %w", err)
	}
	return resp, nil
}

func (t *UserGrpcTransport) FollowUser(ctx context.Context, followerID string, followingID string) (*pb.FollowUserResponse, error) {
	resp, err := t.client.FollowUser(ctx, &pb.FollowUserRequest{
		FollowerId:  followerID,
		FollowingId: followingID,
	})
	if err != nil {
		return nil, fmt.Errorf("cannot FollowUser: %w", err)
	}
	return resp, nil
}

func (t *UserGrpcTransport) UnfollowUser(ctx context.Context, followerID string, followingID string) (*pb.UnfollowUserResponse, error) {
	resp, err := t.client.UnfollowUser(ctx, &pb.UnfollowUserRequest{
		FollowerId:  followerID,
		FollowingId: followingID,
	})
	if err != nil {
		return nil, fmt.Errorf("cannot FollowUser: %w", err)
	}
	return resp, nil
}

func (t *UserGrpcTransport) GetAllFollowers(ctx context.Context, userID string) (*[]entity.RegisterUserRequest, error) {
	resp, err := t.client.GetAllFollowers(ctx, &pb.GetAllFollowersRequest{
		UserId: userID,
	})
	if err != nil {
		return nil, fmt.Errorf("cannot GetAllFollowers: %w", err)
	}

	var followers []entity.RegisterUserRequest

	for {
		res, err := resp.Recv()
		if err != nil {
			if err == io.EOF {
				return &followers, nil
			}
			return nil, err
		}
		follower := res.GetFollowers()

		followerID, err := uuid.Parse(follower.GetId())
		if err != nil {
			return nil, fmt.Errorf("converting id to uuid err; %v", err)
		}

		followers = append(followers, entity.RegisterUserRequest{
			ID:          followerID,
			FirstName:   follower.GetFirstName(),
			LastName:    follower.GetLastName(),
			Username:    follower.GetUsername(),
			Email:       follower.GetEmail(),
			Password:    follower.GetPassword(),
			IsConfirmed: follower.GetIsConfirmed(),
		})
	}
}

func (t *UserGrpcTransport) GetAllFollowings(ctx context.Context, userID string) (*[]entity.RegisterUserRequest, error) {
	resp, err := t.client.GetAllFollowings(ctx, &pb.GetAllFollowingsRequest{
		UserId: userID,
	})
	if err != nil {
		return nil, fmt.Errorf("cannot GetAllFollowings: %w", err)
	}
	var followers []entity.RegisterUserRequest

	for {
		res, err := resp.Recv()
		if err != nil {
			if err == io.EOF {
				return &followers, nil
			}
			return nil, err
		}
		following := res.GetFollowings()

		followingID, err := uuid.Parse(following.GetId())
		if err != nil {
			return nil, fmt.Errorf("converting id to uuid err; %v", err)
		}

		followers = append(followers, entity.RegisterUserRequest{
			ID:          followingID,
			FirstName:   following.GetFirstName(),
			LastName:    following.GetLastName(),
			Username:    following.GetUsername(),
			Email:       following.GetEmail(),
			Password:    following.GetPassword(),
			IsConfirmed: following.GetIsConfirmed(),
		})
	}
}
