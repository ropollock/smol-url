package model

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID          uuid.UUID `gorm:"type:uuid;default:gen_random_uuid()" json:"id,omitempty"`
	Name        string    `json:"name,omitempty"`
	Username    string    `gorm:"index" json:"username,omitempty"`
	Password    string    `json:"-"`
	Email       string    `json:"email,omitempty"`
	CreatedTS   time.Time `json:"created_ts"`
	LastLoginTS time.Time `json:"last_login_ts"`
	IsAdmin     bool      `json:"is_admin,omitempty"`
}
