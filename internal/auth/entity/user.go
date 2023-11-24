package entity

import "github.com/google/uuid"

// swagger:model RegisterUserRequest
type RegisterUserRequest struct {
	ID          uuid.UUID `json:"id"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	Username    string    `json:"username"`
	Email       string    `json:"email"`
	Password    string    `json:"password"`
	IsConfirmed bool      `json:"is_confirmed"`
	RoleID      int64     `json:"role_id"`
}
