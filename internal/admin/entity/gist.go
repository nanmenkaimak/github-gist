package entity

import (
	"time"

	"github.com/google/uuid"
)

type Gist struct {
	ID          uuid.UUID `json:"id"`
	UserID      uuid.UUID `json:"user_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Visible     bool      `json:"visible"`
	IsForked    bool      `json:"is_forked"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Commit struct {
	ID        uuid.UUID `json:"id"`
	GistID    uuid.UUID `json:"gist_id"`
	Gist      Gist      `json:"-"`
	Comment   string    `json:"comment"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type File struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	CommitID  uuid.UUID `json:"commit_id"`
	Commit    Commit    `json:"-"`
	Code      string    `json:"code"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// swagger:model GistRequest
type GistRequest struct {
	Gist   Gist   `json:"gist"`
	Commit Commit `json:"commit"`
	Files  []File `json:"files"`
}
