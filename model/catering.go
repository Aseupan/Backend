package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type UserCateringDonation struct {
	ID              uint          `gorm:"primaryKey" json:"id"`
	UserID          uuid.UUID     `json:"user_id"`
	Campaign        Campaign      `gorm:"ForeignKey:CampaignID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CampaignID      uint          `json:"campaign_id"`
	Catering        Catering      `gorm:"ForeignKey:CateringID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CateringID      uint          `json:"catering_id"`
	AdditionalChips pq.Int64Array `gorm:"type:integer[]" json:"additional_chips"`
	PickUp          int           `json:"pickup"`
	IsDone          bool          `gorm:"default:false;not null" json:"is_done"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type CompanyCateringDonation struct {
	ID              uint          `gorm:"primaryKey" json:"id"`
	Company         uuid.UUID     `json:"company_id"`
	Campaign        Campaign      `gorm:"ForeignKey:CampaignID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CampaignID      uint          `json:"campaign_id"`
	Catering        Catering      `gorm:"ForeignKey:CateringID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CateringID      uint          `json:"catering_id"`
	AdditionalChips pq.Int64Array `gorm:"type:integer[]" json:"additional_chips"`
	PickUp          int           `json:"pickup"`
	IsDone          bool          `gorm:"default:false;not null" json:"is_done"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type Catering struct {
	User            User      `gorm:"ForeignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	UserID          uuid.UUID `json:"user_id"`
	CompanyID       uuid.UUID `json:"company_id"`
	ID              uint      `gorm:"primaryKey" json:"id"`
	Name            string    `json:"name"`
	Phone           string    `json:"phone"`
	Address         string    `json:"address"`
	AddressDetailed string    `json:"address_detailed"`
	IsSaved         bool      `json:"is_saved"`
	CreatedAt       time.Time
}

type NewCateringInput struct {
	Name            string `gorm:"binding:required" json:"name"`
	Phone           string `gorm:"binding:required" json:"phone"`
	Address         string `gorm:"binding:required" json:"address"`
	AddressDetailed string `json:"address_detailed"`
	IsSaved         bool   `gorm:"binding:required" json:"is_saved"`
}
