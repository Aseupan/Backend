package Entities

import "time"

type User struct {
	ID        uint      `gorm:"primaryKey"`
	Name      string    `gorm:"not null" json:"name"`
	Username  string    `gorm:"not null" json:"username"`
	Email     string    `gorm:"not null" json:"email"`
	Password  string    `gorm:"not null" json:"password"`
	Points    int       `gorm:"not null" json:"points"`
	CreatedAt time.Time `json:"created_at"`
	EditedAt  time.Time `json:"edited_at"`
}
