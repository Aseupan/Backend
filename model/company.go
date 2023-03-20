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
