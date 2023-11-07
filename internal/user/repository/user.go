package repository

import (
	"github.com/google/uuid"
	"github.com/nanmenkaimak/github-gist/internal/user/entity"
)

func (r *Repo) CreateUser(newUser entity.User) (uuid.UUID, error) {
	err := r.main.Db.Create(&newUser).Error
	if err != nil {
		return uuid.Nil, err
	}
	return newUser.ID, err
}

func (r *Repo) GetUserByUsername(username string) (*entity.User, error) {
	var userByUsername entity.User

	err := r.replica.Db.Where("username = ?", username).Find(&userByUsername).Error
	if err != nil {
		return nil, err
	}
	return &userByUsername, err
}

func (r *Repo) ConfirmUser(email string) error {
	err := r.main.Db.Model(&entity.User{}).Where("email = ?", email).Update("is_confirmed", true).Error
	return err
}

func (r *Repo) GetUserByID(userID string) (*entity.User, error) {
	var userByID entity.User

	err := r.replica.Db.Where("id = ?", userID).Find(&userByID).Error
	if err != nil {
		return nil, err
	}
	return &userByID, err
}
