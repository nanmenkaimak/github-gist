package gist

import (
	"github.com/google/uuid"
	"github.com/nanmenkaimak/github-gist/internal/gist/entity"
)

// swagger:model CreateGistResponse
type CreateGistResponse struct {
	GistID uuid.UUID `json:"gist_id"`
}

type DeleteRequest struct {
	GistID    uuid.UUID
	CommentID uuid.UUID
	UserID    uuid.UUID
	Username  string
}

type GetGistRequest struct {
	GistID     uuid.UUID
	Username   string
	UserID     uuid.UUID
	Visibility bool
	Searching  string
}

type UpdateGistRequest struct {
	Username string
	UserID   uuid.UUID
	Gist     entity.GistRequest
}

type GetAllGistsRequest struct {
	Direction string
	Sort      string
}

type OtherGistRequest struct {
	Username string
}

type ForkRequest struct {
	GistID   uuid.UUID
	UserID   uuid.UUID
	Username string
}

type ForkGistResponse struct {
	GistID uuid.UUID `json:"gist_id"`
}

type UpdateCommentRequest struct {
	CommentID uuid.UUID
	Username  string
	UserID    uuid.UUID
	Text      string
}

type FollowRequest struct {
	FollowerID  uuid.UUID
	FollowingID uuid.UUID
	Username    string
}
