package model

import "github.com/google/uuid"

type UserRequest struct {
	ID       uuid.UUID `param:"id"`
	Name     string    `json:"name"`
	Username string    `json:"username"`
	Password string    `json:"password"`
	Email    string    `json:"email"`
	IsAdmin  bool      `json:"is_admin"`
}
