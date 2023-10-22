package entity

import (
	"github.com/google/uuid"
	"time"
)

type Comment struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid; default:gen_random_uuid()"`
	UserID    uuid.UUID `json:"user_id" gorm:"type:uuid; not null"`
	GistID    uuid.UUID `json:"gist_id" gorm:"not null"`
	Gist      Gist      `json:"-" gorm:"constraint:OnDelete:CASCADE"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
