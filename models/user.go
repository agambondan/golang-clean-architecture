package models

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	UUID        uuid.UUID `sql:"primary_key" json:"id,omitempty"`
	FirstName   string    `json:"first_name,omitempty"`
	LastName    string    `json:"last_name,omitempty"`
	Email       string    `json:"email,omitempty"`
	PhoneNumber string    `json:"phone_number,omitempty"`
	Username    string    `json:"username,omitempty"`
	Password    string    `json:"password,omitempty"`
	RoleId      uint64    `json:"role_id,omitempty"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
	DeletedAt   time.Time `sql:"index" json:"deleted_at,omitempty"`
}