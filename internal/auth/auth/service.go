package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

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
)

type Service struct {
	repo                     repository.Repository
	userGrpcTransport        *transport.UserGrpcTransport
	jwtSecretKey             string
	userVerificationProducer *kafka.Producer
	dbRedis                  *redis.Client
}

func NewAuthService(repo repository.Repository, authConfig config.Auth, userVerificationProducer *kafka.Producer,
	dbRedis *redis.Client, userGrpcTransport *transport.UserGrpcTransport) UseCase {
	return &Service{
		repo:                     repo,
		jwtSecretKey:             authConfig.JwtSecretKey,
		userVerificationProducer: userVerificationProducer,
		dbRedis:                  dbRedis,
		userGrpcTransport:        userGrpcTransport,
	}
}

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

	userToken := entitiy.UserToken{
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

func (a *Service) RegisterUser(ctx context.Context, request entitiy.RegisterUserRequest) (*RegisterUserResponse, error) {
	// user create by grpc
	hashPass, err := a.hashPassword(request.Password)
	if err != nil {
		return nil, fmt.Errorf("hashing password err: %v", err)
	}

	request.Password = hashPass
	response, err := a.userGrpcTransport.CreateUser(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("CreateUser request err: %v", err)
	}

	userID, err := uuid.Parse(response.GetId())
	if err != nil {
		return nil, fmt.Errorf("converting id to uuid err; %v", err)
	}

	resp := &RegisterUserResponse{
		ID: userID,
	}

	randNum1 := rand.Intn(10)
	randNum2 := rand.Intn(10)
	randNum3 := rand.Intn(10)
	randNum4 := rand.Intn(10)

	msg := dto.UserCode{
		Code: fmt.Sprintf("%d%d%d%d", randNum1, randNum2, randNum3, randNum4),
		Key:  request.Email,
	}

	b, err := json.Marshal(&msg)
	if err != nil {
		return resp, fmt.Errorf("failed to marshall UserCode err: %w", err)
	}

	a.userVerificationProducer.ProduceMessage(b)

	return resp, nil
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
	_, err = a.userGrpcTransport.ConfirmUser(ctx, request.Email)
	if err != nil {
		return fmt.Errorf("ConfirmUser request err: %v", err)
	}
	return nil
}

func (a *Service) UpdateUser(ctx context.Context, updatedUser entitiy.RegisterUserRequest) error {
	user, err := a.userGrpcTransport.GetUserByUsername(ctx, updatedUser.Username)
	if err != nil {
		return fmt.Errorf("GetUserByUsername request err: %v", err)
	}
	userID, err := uuid.Parse(user.Id)
	if err != nil {
		return fmt.Errorf("parse uuid err: %v", err)
	}

	if userID != updatedUser.ID {
		return fmt.Errorf("it is not your account err: %v", err)
	}

	_, err = a.userGrpcTransport.UpdateUser(ctx, updatedUser)
	if err != nil {
		return fmt.Errorf("UpdateUser request err: %v", err)
	}
	return nil
}

func (a *Service) ResetCode(ctx context.Context, request ResetCodeRequest) error {
	user, err := a.userGrpcTransport.GetUserByUsername(ctx, request.Username)
	if err != nil {
		return fmt.Errorf("GetUserByUsername request err: %v", err)
	}
	randNum1 := rand.Intn(10)
	randNum2 := rand.Intn(10)
	randNum3 := rand.Intn(10)
	randNum4 := rand.Intn(10)

	msg := &dto.UserCode{
		Code: fmt.Sprintf("%d%d%d%d", randNum1, randNum2, randNum3, randNum4),
		Key:  user.Email,
	}

	b, err := json.Marshal(&msg)
	if err != nil {
		return fmt.Errorf("failed to marshall UserCodeReset err: %w", err)
	}

	fmt.Println(msg)

	a.userVerificationProducer.ProduceMessage(b)

	return nil
}

func (a *Service) ResetPassword(ctx context.Context, request UpdatePasswordRequest) error {
	user, err := a.userGrpcTransport.GetUserByUsername(ctx, request.Username)
	if err != nil {
		return fmt.Errorf("GetUserByUsername request err: %v", err)
	}
	res, err := a.dbRedis.Get(ctx, user.Email).Result()
	if err != nil {
		return fmt.Errorf("redis get err: %v", err)
	}

	if res != request.Code {
		return fmt.Errorf("wrong confirm code")
	}

	hashPass, err := a.hashPassword(request.NewPassword)
	if err != nil {
		return fmt.Errorf("hashing password err: %v", err)
	}

	_, err = a.userGrpcTransport.UpdatePassword(ctx, user.Email, hashPass)
	if err != nil {
		return fmt.Errorf("UpdatePassword request err: %v", err)
	}

	return err
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

type ContextUserKey struct{}

type ContextUser struct {
	ID uuid.UUID `json:"user_id"`
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
