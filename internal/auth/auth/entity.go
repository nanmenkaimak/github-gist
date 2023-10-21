package auth

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type GenerateTokenRequest struct {
	Username string
	Password string
}

type JwtUserToken struct {
	Token        string
	RefreshToken string
}

type JwtRenewToken struct {
	Token string
}

type MyCustomClaims struct {
	UserId uuid.UUID `json:"user_id"`
	jwt.RegisteredClaims
}
