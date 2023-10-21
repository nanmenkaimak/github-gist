package auth

import "context"

type UseCase interface {
	GenerateToken(ctx context.Context, request GenerateTokenRequest) (*JwtUserToken, error)
	RenewToken(ctx context.Context, refreshToken string) (*JwtRenewToken, error)
}
