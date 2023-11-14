package auth

import (
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/nanmenkaimak/github-gist/internal/gist/config"
)

type ContextUserKey struct{}

type ContextUser struct {
	ID uuid.UUID `json:"user_id"`
}

type Service struct {
	jwtSecretKey string
}

func NewService(authConfig config.Auth) *Service {
	return &Service{
		jwtSecretKey: authConfig.JwtSecretKey,
	}
}

func (a *Service) GetJWTUser(jwtToken string) (*ContextUser, error) {
	token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("wrong token")
		}
		return []byte(a.jwtSecretKey), nil
	})

	if err != nil {
		return nil, fmt.Errorf("token parsing err: %v", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("unexpected type %T", claims)
	}

	user, err := a.getUserFromJWT(claims)
	if err != nil {
		return nil, fmt.Errorf("failed to get user from jwt err: %w", err)
	}
	return user, nil
}

func (a *Service) getUserFromJWT(claims jwt.MapClaims) (*ContextUser, error) {
	user := &ContextUser{}
	userID, err := uuid.Parse(claims["user_id"].(string))
	if err != nil {
		return nil, fmt.Errorf("formating user_id into uuid err: %v", err)
	}

	user.ID = userID

	return user, nil
}

func (a *Service) GetContextUserKey() ContextUserKey {
	return ContextUserKey{}
}
