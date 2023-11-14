package entitiy

import (
	"time"

	"github.com/google/uuid"
)

type UserToken struct {
	ID           uuid.UUID `json:"id" gorm:"primaryKey; type:uuid; default:gen_random_uuid()"`
	UserID       uuid.UUID `json:"user_id" gorm:"type:uuid; not null"`
	Token        string    `json:"token" gorm:"not null"`
	RefreshToken string    `json:"refresh_token" gorm:"not null"`
	ExpiresAt    time.Time `json:"expires_at" gorm:"default:now()"`
	CreatedAt    time.Time `json:"created_at" gorm:"default:now()"`
}
