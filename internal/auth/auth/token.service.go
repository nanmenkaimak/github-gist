package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/nanmenkaimak/github-gist/internal/auth/entity"
	"golang.org/x/crypto/bcrypt"
)

func (a *Service) GenerateToken(ctx context.Context, request GenerateTokenRequest) (*JwtUserToken, error) {
	user, err := a.userGrpcTransport.GetUserByUsername(ctx, request.Username)
	if err != nil {
		return nil, fmt.Errorf("GetUserByUsername request err: %v", err)
	}

	err = a.comparePassword(user.Password, request.Password)
	if err != nil {
		return nil, fmt.Errorf("comparePassword err: %v", err)
	}

	userID, err := uuid.Parse(user.Id)
	if err != nil {
		return nil, fmt.Errorf("parse uuid err: %v", err)
	}

	claims := MyCustomClaims{
		user.RoleId,
		userID,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	secretKey := []byte(a.jwtSecretKey)
	claimToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := claimToken.SignedString(secretKey)
	if err != nil {
		return nil, fmt.Errorf("SignedString err: %w", err)
	}

	rClaims := MyCustomClaims{
		user.RoleId,
		userID,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(40 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	rClaimToken := jwt.NewWithClaims(jwt.SigningMethodHS256, rClaims)

	refreshTokenString, err := rClaimToken.SignedString(secretKey)
	if err != nil {
		return nil, fmt.Errorf("SignedString err: %w", err)
	}

	userToken := entity.UserToken{
		UserID:       userID,
		Token:        tokenString,
		RefreshToken: refreshTokenString,
	}

	err = a.repo.CreateUserToken(userToken)
	if err != nil {
		return nil, fmt.Errorf("CreateUserToken err: %w", err)
	}

	jwtToken := &JwtUserToken{
		Token:        tokenString,
		RefreshToken: refreshTokenString,
	}

	return jwtToken, nil
}

func (a *Service) RenewToken(ctx context.Context, refreshToken string) (*JwtRenewToken, error) {
	claims, err := a.parseToken(refreshToken)
	if err != nil {
		return nil, fmt.Errorf("parse refresh token err: %v", err)
	}

	userID, err := uuid.Parse(claims["user_id"].(string))
	if err != nil {
		return nil, fmt.Errorf("convert refresh token err: %v", err)
	}

	roleID := int64(claims["role_id"].(float64))

	newClaims := MyCustomClaims{
		roleID,
		userID,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	newClaimsToken := jwt.NewWithClaims(jwt.SigningMethodHS256, newClaims)
	tokenString, err := newClaimsToken.SignedString([]byte(a.jwtSecretKey))
	if err != nil {
		return nil, fmt.Errorf("SignedString err: %w", err)
	}

	newToken := entity.UserToken{
		UserID: userID,
		Token:  tokenString,
	}
	err = a.repo.UpdateUserToken(newToken)
	if err != nil {
		return nil, fmt.Errorf("UpdateUserToken err: %w", err)
	}
	jwtToken := &JwtRenewToken{
		Token: tokenString,
	}

	return jwtToken, nil
}

type ContextUserKey struct{}

type ContextUser struct {
	ID     uuid.UUID `json:"user_id"`
	RoleID int64     `json:"role_id"`
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

	roleID := int64(claims["role_id"].(float64))

	user.ID = userID
	user.RoleID = roleID

	return user, nil
}

func (a *Service) parseToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("wrong token")
		}
		return []byte(a.jwtSecretKey), nil
	})

	if err != nil {
		return nil, fmt.Errorf("token parsing err: %v", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("token is not valid")
}

func (a *Service) comparePassword(password1 string, password2 string) error {
	err := bcrypt.CompareHashAndPassword([]byte(password1), []byte(password2))
	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return fmt.Errorf("incorrect password err: %v", err)
	} else if err != nil {
		return fmt.Errorf("password auth err: %v", err)
	}
	return nil
}
