package entity

import (
	"time"

	"github.com/google/uuid"
)

// swagger:model User
type User struct {
	ID          uuid.UUID `json:"id"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	Username    string    `json:"username"`
	Email       string    `json:"email"`
	Password    string    `json:"password"`
	IsConfirmed bool      `json:"is_confirmed"`
	RoleID      int       `json:"role_id"`
	Role        Role      `json:"role"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Role struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
