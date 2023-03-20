package model

import (
	"time"

	"github.com/google/uuid"
)

type Company struct {
	ID             uuid.UUID `gorm:"primaryKey" json:"id"`
	CompanyName    string    `json:"company_name"`
	CompanyAddress string    `json:"company_address"`
	CompanyEmail   string    `json:"company_email"`
	Password       string    `json:"password"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type CompanyRegisterInput struct {
	CompanyName     string `gorm:"binding:required" json:"company_name"`
	CompanyAddress  string `gorm:"binding:required" json:"company_address"`
	CompanyEmail    string `gorm:"binding:required" json:"company_email"`
	Password        string `gorm:"binding:required" json:"password"`
	ConfirmPassword string `gorm:"binding:required" json:"confirm_password"`
}
