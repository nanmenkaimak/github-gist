package gist

import "github.com/google/uuid"

type CreateGistResponse struct {
	GistID uuid.UUID `json:"gist_id"`
}

type GetGistByIDRequest struct {
	GistID uuid.UUID
}
