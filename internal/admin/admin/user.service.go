package admin

import (
	"context"
	"fmt"
	"github.com/nanmenkaimak/github-gist/internal/admin/entity"
	"golang.org/x/crypto/bcrypt"
)

func (a *Service) UpdateUserByUsername(ctx context.Context, request UpdateUserRequest) error {
	ok, err := a.userRepo.IsAdmin(request.UserID)
	if err != nil {
		return fmt.Errorf("IsAdmin request err: %v", err)
	}
	if !ok {
		return fmt.Errorf("you are not admin")
	}
	hashPass, err := a.hashPassword(request.UpdatedUser.Password)
	if err != nil {
		return fmt.Errorf("hashing password err: %v", err)
	}

	request.UpdatedUser.Password = hashPass

	err = a.userRepo.UpdateUser(request.UpdatedUser)
	if err != nil {
		return fmt.Errorf("UpdateUser request err: %v", err)
	}
	return nil
}

func (a *Service) GetAllUsers(ctx context.Context, request GetUserRequest) (*[]entity.User, error) {
	ok, err := a.userRepo.IsAdmin(request.UserID)
	if err != nil {
		return nil, fmt.Errorf("IsAdmin request err: %v", err)
	}
	if !ok {
		return nil, fmt.Errorf("you are not admin")
	}
	users, err := a.userRepo.GetAllUsers()
	if err != nil {
		return nil, fmt.Errorf("GetAllUsers request err: %v", err)
	}
	return users, nil
}

func (a *Service) GetUserByUsername(ctx context.Context, request GetUserRequest) (*entity.User, error) {
	ok, err := a.userRepo.IsAdmin(request.UserID)
	if err != nil {
		return nil, fmt.Errorf("IsAdmin request err: %v", err)
	}
	if !ok {
		return nil, fmt.Errorf("you are not admin")
	}
	user, err := a.userRepo.GetUserByUsername(request.UpdatedUserUsername)
	if err != nil {
		return nil, fmt.Errorf("GetUserByUsername request err: %v", err)
	}
	return user, nil
}

func (a *Service) hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}
