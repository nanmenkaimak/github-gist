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

type ConfirmUserRequest struct {
	Email string
	Code  string
}

type RegisterUserResponse struct {
	ID uuid.UUID
}

type ResetCodeRequest struct {
	Username string
	UserID   uuid.UUID
}

type UpdatePasswordRequest struct {
	Username    string `json:"username"`
	NewPassword string `json:"new_password"`
	Code        string `json:"code"`
}

type FollowRequest struct {
	FollowerID  uuid.UUID
	FollowingID uuid.UUID
	Username    string
}
