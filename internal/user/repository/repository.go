package repository

import (
	"github.com/google/uuid"
	"github.com/nanmenkaimak/github-gist/internal/user/database/dbpostgres"
	"github.com/nanmenkaimak/github-gist/internal/user/entity"
)

type Repository interface {
	UserRepository
	FollowRepository
}

type UserRepository interface {
	CreateUser(newUser entity.User) (uuid.UUID, error)
	GetUserByUsername(username string) (*entity.User, error)
	ConfirmUser(email string) error
	GetUserByID(userID string) (*entity.User, error)
	UpdateUser(updatedUser entity.User) error
}

type FollowRepository interface {
	FollowUser(follower entity.Follower) error
	UnfollowUser(follower entity.Follower) error
	GetAllFollowers(userID string) ([]entity.User, error)
	GetAllFollowings(userID string) ([]entity.User, error)
}

type Repo struct {
	main    *dbpostgres.Db
	replica *dbpostgres.Db
}

func NewRepository(connMain *dbpostgres.Db, connReplica *dbpostgres.Db) *Repo {
	return &Repo{
		main:    connMain,
		replica: connReplica,
	}
}
