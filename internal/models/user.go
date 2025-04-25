package models

import (
	"time"

	"github.com/google/uuid"
)

type Role string

const (
	RoleUser  Role = "User"
	RoleAdmin Role = "Admin"
)

type User struct {
	ID              uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Email           string    `gorm:"unique;not null"`
	PasswordHash    string    `gorm:"not null"`
	Name            string
	Phone           string
	AvatarURL       string
	Role            Role `gorm:"type:role_enum;default:'User'"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
	LastLogin       time.Time
	IsEmailVerified bool
}
