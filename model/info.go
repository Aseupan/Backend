package model

import "time"

type Info struct {
	ID        uint   `json:"id" gorm:"primaryKey"`
	Type      string `json:"type"`
	Title     string `json:"title"`
	Body      string `json:"body"`
	Thumbnail string `json:"thumbnail"`
	CreatedAt time.Time
}
