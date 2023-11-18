package entity

import (
	"time"

	"github.com/google/uuid"
)

// swagger:model GistRequest
type GistRequest struct {
	Gist   Gist   `json:"gist"`
	Commit Commit `json:"commit"`
	Files  []File `json:"files"`
}

type Gist struct {
	ID          uuid.UUID `json:"id" gorm:"primaryKey; type:uuid; default:gen_random_uuid()"`
	UserID      uuid.UUID `json:"user_id" gorm:"type:uuid; not null"`
	Name        string    `json:"name" gorm:"unique; type:varchar(50); not null"`
	Description string    `json:"description" gorm:"type:varchar(150)"`
	Visible     bool      `json:"visible" gorm:"default:false"`
	IsForked    bool      `json:"is_forked" gorm:"default:false"`
	CreatedAt   time.Time `json:"created_at" gorm:"default:now()"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"default:now()"`
}

type Commit struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid; default:gen_random_uuid()"`
	GistID    uuid.UUID `json:"gist_id" gorm:"not null"`
	Gist      Gist      `json:"-" gorm:"constraint:OnDelete:CASCADE"`
	Comment   string    `json:"comment" gorm:"type:varchar(250); not null"`
	CreatedAt time.Time `json:"created_at" gorm:"default:now()"`
	UpdatedAt time.Time `json:"updated_at" gorm:"default:now()"`
}

type File struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid; default:gen_random_uuid()"`
	Name      string    `json:"name" gorm:"type:varchar(50); not null"`
	CommitID  uuid.UUID `json:"commit_id" gorm:"not null"`
	Commit    Commit    `json:"-" gorm:"constraint:OnDelete:CASCADE"`
	Code      string    `json:"code" gorm:"not null"`
	CreatedAt time.Time `json:"created_at" gorm:"default:now()"`
	UpdatedAt time.Time `json:"updated_at" gorm:"default:now()"`
}
