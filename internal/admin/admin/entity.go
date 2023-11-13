package admin

import (
	"github.com/google/uuid"
	"github.com/nanmenkaimak/github-gist/internal/admin/entity"
)

type GetAllGistsRequest struct {
	Direction string
	Sort      string
	UserID    uuid.UUID
}

type GetGistRequest struct {
	GistID uuid.UUID
	UserID uuid.UUID
}

type UpdateUserRequest struct {
	UpdatedUser entity.User
	UserID      uuid.UUID
}

type GetUserRequest struct {
	UserID              uuid.UUID
	UpdatedUserUsername string
}
