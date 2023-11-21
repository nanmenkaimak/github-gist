package entity

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID          uuid.UUID `json:"id" gorm:"primaryKey; type:uuid; default:gen_random_uuid()"`
	FirstName   string    `json:"first_name" gorm:"type:varchar(50); not null"`
	LastName    string    `json:"last_name" gorm:"type:varchar(50); not null"`
	Username    string    `json:"username" gorm:"unique; type:varchar(50); not null"`
	Email       string    `json:"email" gorm:"unique; type:varchar(100); not null"`
	Password    string    `json:"password" gorm:"type:varchar(255); not null"`
	IsConfirmed bool      `json:"is_confirmed" gorm:"default:false"`
	RoleID      int       `json:"role_id" gorm:"default:1"`
	Role        Role      `json:"-"`
	CreatedAt   time.Time `json:"created_at" gorm:"default:now()"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"default:now()"`
}

type Role struct {
	ID   int    `json:"id"`
	Name string `json:"name" gorm:"not null"`
}

type Follower struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid; default:gen_random_uuid()"`
	FollowingID uuid.UUID `json:"following_id" gorm:"not null"`
	Following   User      `json:"-" gorm:"constraint:OnDelete:CASCADE"`
	FollowerID  uuid.UUID `json:"follower_id" gorm:"not null"`
	Follower    User      `json:"-" gorm:"constraint:OnDelete:CASCADE"`
	CreatedAt   time.Time `json:"created_at" gorm:"default:now()"`
}
