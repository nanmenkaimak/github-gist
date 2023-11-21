package entity

import "github.com/google/uuid"

type RegisterUserRequest struct {
	ID          uuid.UUID `json:"id"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	Username    string    `json:"username"`
	Email       string    `json:"email"`
	Password    string    `json:"password"`
	IsConfirmed bool      `json:"is_confirmed"`
}
