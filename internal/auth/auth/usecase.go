package auth

import (
	"context"

	"github.com/nanmenkaimak/github-gist/internal/auth/entitiy"
)

type UseCase interface {
	GenerateToken(ctx context.Context, request GenerateTokenRequest) (*JwtUserToken, error)
	RenewToken(ctx context.Context, refreshToken string) (*JwtRenewToken, error)
	RegisterUser(ctx context.Context, request entitiy.RegisterUserRequest) (*RegisterUserResponse, error)
	ConfirmUser(ctx context.Context, request ConfirmUserRequest) error
	UpdateUser(ctx context.Context, updatedUser entitiy.RegisterUserRequest) error
	GetJWTUser(jwtToken string) (*ContextUser, error)
	ResetCode(ctx context.Context, request ResetCodeRequest) error
	ResetPassword(ctx context.Context, request UpdatePasswordRequest) error
}
