package model

import (
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type Campaign struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	CompanyID   uuid.UUID      `json:"company_id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Progress    int            `json:"progress"`
	Target      int            `json:"target"`
	Area        string         `json:"area"`
	StartDate   string         `json:"start_date"`
	EndDate     string         `json:"end_date"`
	Thumbnail1  string         `json:"thumbnail_1"`
	Thumbnail2  string         `json:"thumbnail_2"`
	Thumbnail3  string         `json:"thumbnail_3"`
	Thumbnail4  string         `json:"thumbnail_4"`
	Thumbnail5  string         `json:"thumbnail_5"`
	Urgent      int            `json:"urgent"`
	Type        pq.StringArray `gorm:"type:text[]" json:"type"`
}

type CampaignInput struct {
	Name        string         `gorm:"binding:required" json:"name"`
	Description string         `gorm:"binding:required" json:"description"`
	Target      int            `gorm:"binding:required" json:"target"`
	Area        string         `gorm:"binding:required" json:"area"`
	StartDate   string         `gorm:"binding:required" json:"start_date"`
	EndDate     string         `gorm:"binding:required" json:"end_date"`
	Thumbnail   string         `gorm:"binding:required" json:"thumbnail"`
	Urgent      int            `gorm:"binding:required" json:"urgent"`
	Type        pq.StringArray `gorm:"type:text[]" json:"type"`
}
