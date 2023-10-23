package gist

import "github.com/google/uuid"

type CreateGistResponse struct {
	GistID uuid.UUID `json:"gist_id"`
}

type GetGistRequest struct {
	GistID   uuid.UUID
	Username string
	UserID   uuid.UUID
}
