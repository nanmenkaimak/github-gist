package transport

import (
	"context"
	"fmt"
	"github.com/nanmenkaimak/github-gist/internal/gist/config"
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

func (t *UserGrpcTransport) GetUserByID(ctx context.Context, userID string) (*pb.User, error) {
	resp, err := t.client.GetUserByID(ctx, &pb.GetUserByIDRequest{
		Id: userID,
	})
	if err != nil {
		return nil, fmt.Errorf("cannot GetUserByID: %w", err)
	}

	if resp == nil {
		return nil, fmt.Errorf("not found")
	}

	return resp.User, nil
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
