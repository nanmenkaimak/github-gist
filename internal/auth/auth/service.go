package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/nanmenkaimak/github-gist/internal/auth/config"
	"github.com/nanmenkaimak/github-gist/internal/auth/controller/consumer/dto"
	"github.com/nanmenkaimak/github-gist/internal/auth/entitiy"
	"github.com/nanmenkaimak/github-gist/internal/auth/repository"
	"github.com/nanmenkaimak/github-gist/internal/auth/transport"
	"github.com/nanmenkaimak/github-gist/internal/kafka"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"time"
)

type Service struct {
	repo                     repository.Repository
	userTransport            *transport.UserTransport
	jwtSecretKey             string
	userVerificationProducer *kafka.Producer
	dbRedis                  *redis.Client
}

func NewAuthService(repo repository.Repository, authConfig config.Auth,
	userTransport *transport.UserTransport, userVerificationProducer *kafka.Producer, dbRedis *redis.Client) UseCase {
	return &Service{
		repo:                     repo,
		userTransport:            userTransport,
		jwtSecretKey:             authConfig.JwtSecretKey,
		userVerificationProducer: userVerificationProducer,
		dbRedis:                  dbRedis,
	}
}

func (a *Service) GenerateToken(ctx context.Context, request GenerateTokenRequest) (*JwtUserToken, error) {
	user, err := a.userTransport.GetUser(ctx, request.Username)
	if err != nil {
		return nil, fmt.Errorf("GetUser request err: %v", err)
	}

	err = a.comparePassword(user.Password, request.Password)
	if err != nil {
		return nil, fmt.Errorf("comparePassword err: %v", err)
	}

	claims := MyCustomClaims{
		user.ID,
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
		user.ID,
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

	userToken := entitiy.UserToken{
		UserID:       user.ID,
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

	newClaims := MyCustomClaims{
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

	newToken := entitiy.UserToken{
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

func (a *Service) RegisterUser(ctx context.Context) error {
	// user create by grpc

	randNum1 := rand.Intn(10)
	randNum2 := rand.Intn(10)
	randNum3 := rand.Intn(10)
	randNum4 := rand.Intn(10)

	msg := dto.UserCode{
		Code:  fmt.Sprintf("%d%d%d%d", randNum1, randNum2, randNum3, randNum4),
		Email: "aristanovali618@gmail.com",
	}

	b, err := json.Marshal(&msg)
	if err != nil {
		return fmt.Errorf("failed to marshall UserCode err: %w", err)
	}

	a.userVerificationProducer.ProduceMessage(b)

	return nil
}

func (a *Service) ConfirmUser(ctx context.Context, request ConfirmUserRequest) error {
	// check code in database
	res, err := a.dbRedis.Get(ctx, request.Email).Result()
	if err != nil {
		return fmt.Errorf("redis get err: %v", err)
	}

	if res != request.Code {
		return fmt.Errorf("wrong confirm code")
	}
	// if ok update user confirmed by grpc
	return nil
}

func (a *Service) hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
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
