package admin

import (
	"github.com/nanmenkaimak/github-gist/internal/admin/repository"
	"github.com/nanmenkaimak/github-gist/internal/admin/storage"
)

type Service struct {
	gistRepo repository.Gist
	userRepo repository.User
	storage  *storage.DataStorage
}

func NewAdminService(gistRepo repository.Gist, userRepo repository.User, storage *storage.DataStorage) UseCase {
	return &Service{
		gistRepo: gistRepo,
		userRepo: userRepo,
		storage:  storage,
	}
}
