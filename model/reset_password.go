package model

import (
	"time"

	"github.com/google/uuid"
)

type ResetPassword struct {
	ID        uuid.UUID `gorm:"primaryKey" json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	CompanyID uuid.UUID `json:"company_id"`
	Email     string    `gorm:"not null;unique" json:"email"`
	Code      string    `gorm:"not null" json:"code"`
	IsUsed    bool      `gorm:"default:false" json:"is_used"`
	CreatedAt time.Time
}

type ForgotPasswordEmailInput struct {
	Email string `gorm:"binding:required" json:"email"`
}

type ForgotPasswordLogin struct {
	Email string `gorm:"binding:required" json:"email"`
	Code  string `gorm:"binding:required" json:"code"`
}
