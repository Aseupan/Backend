package model

import (
	"time"

	"github.com/google/uuid"
)

type TransactionHistory struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uuid.UUID `gorm:"null" json:"user_id"`
	CompanyID uuid.UUID `gorm:"null" json:"company_id"`
	OrderID   string    `json:"order_id"`
	Price     int       `json:"price"`
	Points    int       `json:"points"`
	CreatedAt time.Time
}
