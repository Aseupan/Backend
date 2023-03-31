package model

import (
	"time"

	"github.com/google/uuid"
)

type Company struct {
	ID             uuid.UUID `gorm:"primaryKey" json:"id"`
	CompanyName    string    `json:"company_name"`
	CompanyPicture string    `json:"company_picture"`
	CompanyAddress string    `json:"company_address"`
	CompanyEmail   string    `json:"company_email"`
	CompanyPhone   string    `json:"company_phone"`
	Password       string    `json:"password"`
	Point          int       `gorm:"default:0;not null" json:"point"`
	Verified       bool      `gorm:"default:false;not null" json:"verified"`
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

type CompanyUpdateProfileInput struct {
	CompanyPicture string `json:"company_picture"`
	CompanyName    string `json:"company_name"`
	CompanyEmail   string `json:"company_email"`
	CompanyPhone   string `json:"company_phone"`
}

type CompanyResetPasswordInput struct {
	Password string `gorm:"binding:required" json:"password"`
}

type CampaignCompanyReciever struct {
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Address string `json:"address"`
}
