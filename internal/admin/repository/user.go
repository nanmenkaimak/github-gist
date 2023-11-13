package repository

import "github.com/nanmenkaimak/github-gist/internal/admin/entity"

func (r *UserRepo) UpdateUser(updatedUser entity.User) error {
	err := r.main.Db.Model(&updatedUser).Where("username = ?", updatedUser.Username).
		Updates(updatedUser).Error
	return err
}

func (r *UserRepo) GetUserByUsername(username string) (*entity.User, error) {
	var userByUsername entity.User

	err := r.replica.Db.Where("username = ?", username).Find(&userByUsername).Error
	if err != nil {
		return nil, err
	}
	return &userByUsername, err
}

func (r *UserRepo) GetAllUsers() (*[]entity.User, error) {
	var allUsers []entity.User

	err := r.replica.Db.Find(&allUsers).Order("created_at desc").Error
	if err != nil {
		return nil, err
	}

	return &allUsers, err
}
