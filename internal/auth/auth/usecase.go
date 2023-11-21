package auth

import (
	"context"

	"github.com/nanmenkaimak/github-gist/internal/auth/entity"
)

type UseCase interface {
	Token
	Follow
	User
}

type Token interface {
	GenerateToken(ctx context.Context, request GenerateTokenRequest) (*JwtUserToken, error)
	RenewToken(ctx context.Context, refreshToken string) (*JwtRenewToken, error)
	GetJWTUser(jwtToken string) (*ContextUser, error)
}

type Follow interface {
	FollowUser(ctx context.Context, request FollowRequest) error
	UnfollowUser(ctx context.Context, request FollowRequest) error
	GetAllFollowers(ctx context.Context, username string) (*[]entity.RegisterUserRequest, error)
	GetAllFollowings(ctx context.Context, username string) (*[]entity.RegisterUserRequest, error)
}

type User interface {
	GetUserInfo(ctx context.Context, username string) (entity.RegisterUserRequest, error)
	RegisterUser(ctx context.Context, request entity.RegisterUserRequest) (*RegisterUserResponse, error)
	ConfirmUser(ctx context.Context, request ConfirmUserRequest) error
	UpdateUser(ctx context.Context, updatedUser entity.RegisterUserRequest) error
	ResetCode(ctx context.Context, request ResetCodeRequest) error
	ResetPassword(ctx context.Context, request UpdatePasswordRequest) error
}
