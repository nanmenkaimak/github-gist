package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/nanmenkaimak/github-gist/internal/auth/controller/consumer/dto"
	"github.com/nanmenkaimak/github-gist/internal/auth/entitiy"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
)

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
