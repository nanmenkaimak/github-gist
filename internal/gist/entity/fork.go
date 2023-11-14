package entity

import (
	"time"

	"github.com/google/uuid"
)

type Fork struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid; default:gen_random_uuid()"`
	GistID    uuid.UUID `json:"gist_id" gorm:"not null"`
	Gist      Gist      `json:"-" gorm:"constraint:OnDelete:CASCADE"`
	NewGistID uuid.UUID `json:"new_gist_id" gorm:"not null"`
	NewGist   Gist      `json:"-" gorm:"constraint:OnDelete:CASCADE"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
