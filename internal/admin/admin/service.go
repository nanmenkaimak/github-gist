package admin

import (
	"github.com/nanmenkaimak/github-gist/internal/admin/repository"
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
