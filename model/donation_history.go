package model

import (
	"time"

	"github.com/google/uuid"
)

type History struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	CompanyID uuid.UUID `json:"company_id"`
	Title     string    `json:"title"`
	Category  int       `json:"category"` // 1 -> donation to campaign, 2 -> buy chips, 3 -> buy reward
	CreatedAt time.Time
}
