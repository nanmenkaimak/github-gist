package auth

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// swagger:model GenerateTokenRequest
type GenerateTokenRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// swagger:model GenerateTokenResponse
type GenerateTokenResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

// swagger:model RenewTokenRequest
type RenewTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

// swagger:model RenewTokenResponse
type RenewTokenResponse struct {
	Token string `json:"token"`
}

type JwtUserToken struct {
	Token        string
	RefreshToken string
}

type JwtRenewToken struct {
	Token string
}

type MyCustomClaims struct {
	RoleID int64     `json:"role_id"`
	UserId uuid.UUID `json:"user_id"`
	jwt.RegisteredClaims
}

// swagger:model ConfirmUserRequest
type ConfirmUserRequest struct {
	Email string
	Code  string
}

// swagger:model RegisterUserResponse
type RegisterUserResponse struct {
	ID uuid.UUID
}

// swagger:model ResetCodeRequest
type ResetCodeRequest struct {
	Username string
	UserID   uuid.UUID
}

// swagger:model UpdatePasswordRequest
type UpdatePasswordRequest struct {
	Username    string `json:"-"`
	NewPassword string `json:"new_password"`
	Code        string `json:"code"`
}

type FollowRequest struct {
	FollowerID  uuid.UUID
	FollowingID uuid.UUID
	Username    string
}
