package model

import (
	"time"

	"github.com/google/uuid"
)

type TransactionHistory struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	User      User      `gorm:"ForeignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	UserID    uuid.UUID `json:"user_id"`
	OrderID   string    `json:"order_id"`
	Price     int       `json:"price"`
	Points    int       `json:"points"`
	CreatedAt time.Time
}
