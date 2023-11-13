package admin

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/nanmenkaimak/github-gist/internal/admin/entity"
	"github.com/nanmenkaimak/github-gist/internal/admin/repository"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	gistRepo repository.Gist
	userRepo repository.User
}

func NewAdminService(gistRepo repository.Gist, userRepo repository.User) UseCase {
	return &Service{
		gistRepo: gistRepo,
		userRepo: userRepo,
	}
}

func (a *Service) GetAllGists(ctx context.Context, request GetAllGistsRequest) (*[]entity.GistRequest, error) {
	ok, err := a.userRepo.IsAdmin(request.UserID)
	if err != nil {
		return nil, fmt.Errorf("IsAdmin request err: %v", err)
	}
	if !ok {
		return nil, fmt.Errorf("you are not admin")
	}

	if request.Sort == "" {
		request.Sort = "created_at"
	}
	if request.Direction == "" {
		request.Direction = "desc"
	}
	allGists, err := a.gistRepo.GetOtherAllGists(request.Sort, request.Direction)
	if err != nil {
		return nil, fmt.Errorf("getting all gists err: %v", err)
	}

	return allGists, nil
}

func (a *Service) GetGistByID(ctx context.Context, request GetGistRequest) (*entity.GistRequest, error) {
	ok, err := a.userRepo.IsAdmin(request.UserID)
	if err != nil {
		return nil, fmt.Errorf("IsAdmin request err: %v", err)
	}
	if !ok {
		return nil, fmt.Errorf("you are not admin")
	}
	gist, err := a.gistRepo.GetGistByID(request.GistID)
	if err != nil {
		return nil, fmt.Errorf("getting gist err: %v", err)
	}

	if gist.Gist.ID == uuid.Nil {
		return nil, nil
	}

	return gist, nil
}

func (a *Service) DeleteGistByID(ctx context.Context, request GetGistRequest) error {
	ok, err := a.userRepo.IsAdmin(request.UserID)
	if err != nil {
		return fmt.Errorf("IsAdmin request err: %v", err)
	}
	if !ok {
		return fmt.Errorf("you are not admin")
	}
	err = a.gistRepo.DeleteGistByID(request.GistID)
	if err != nil {
		return fmt.Errorf("delete gist err: %v", err)
	}
	return nil
}

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

func (a *Service) comparePassword(password1 string, password2 string) error {
	err := bcrypt.CompareHashAndPassword([]byte(password1), []byte(password2))
	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return fmt.Errorf("incorrect password err: %v", err)
	} else if err != nil {
		return fmt.Errorf("password auth err: %v", err)
	}
	return nil
}
