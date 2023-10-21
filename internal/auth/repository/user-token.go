package repository

import (
	"fmt"
	"github.com/nanmenkaimak/github-gist/internal/auth/entitiy"
)

func (r *Repo) CreateUserToken(userToken entitiy.UserToken) error {
	tx := r.main.Db.Create(&userToken)
	if tx.Error != nil {
		return fmt.Errorf("failed creating user token err: %v", tx.Error)
	}
	return nil
}

func (r *Repo) UpdateUserToken(userToken entitiy.UserToken) error {
	tx := r.main.Db.Model(userToken).Where("user_id = ?", userToken.UserID).Update("token", userToken.Token)
	if tx.Error != nil {
		return fmt.Errorf("failed updating user token err: %v", tx.Error)
	}
	return nil
}
