package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID             uuid.UUID `gorm:"primaryKey" json:"id"`
	Name           string    `gorm:"not null" json:"name"`
	ProfilePicture string    `json:"profile_picture"`
	Phone          string    `gorm:"unique;default:null" json:"phone"` // default is null
	Email          string    `gorm:"unique;not null" json:"email"`
	Password       string    `gorm:"not null" json:"password"`
	Point          int       `gorm:"default:0;not null" json:"point"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type LoginInput struct {
	Email    string `gorm:"binding:required" json:"email"`
	Password string `gorm:"binding:required" json:"password"`
}

type UserRegisterInput struct {
	Name            string `gorm:"binding:required" json:"name"`
	Email           string `gorm:"binding:required" json:"email"`
	Password        string `gorm:"binding:required" json:"password"`
	ConfirmPassword string `gorm:"binding:required" json:"confirm_password"`
}

type UserUpdateProfileInput struct {
	ProfilePicture string `json:"profile_picture"`
	Name           string `json:"name"`
	Email          string `json:"email"`
	Phone          string `json:"phone"`
}
type UserResetPasswordInput struct {
	Password string `gorm:"binding:required" json:"password"`
}
