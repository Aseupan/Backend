package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID             uuid.UUID `gorm:"primaryKey" json:"id"`
	Name           string    `json:"name"`
	Phone          string    `gorm:"unique" json:"phone"`
	Email          string    `gorm:"unique" json:"email"`
	Password       string    `json:"password"`
	ProfilePicture string    `json:"profile_picture"`
	Point          int       `json:"point"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type UserLoginInput struct {
	Email    string `gorm:"binding:required" json:"email"`
	Password string `gorm:"binding:required" json:"password"`
}

type UserRegisterInput struct {
	Name            string `gorm:"binding:required" json:"name"`
	Email           string `gorm:"binding:required" json:"email"`
	Password        string `gorm:"binding:required" json:"password"`
	ConfirmPassword string `gorm:"binding:required" json:"confirm_password"`
}

type UserResetPasswordInput struct {
	Password string `gorm:"binding:required" json:"password"`
}
