package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type UserPersonalDonation struct {
	ID              uint          `gorm:"primaryKey" json:"id"`
	UserID          uuid.UUID     `json:"user_id"`
	Campaign        Campaign      `gorm:"ForeignKey:CampaignID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CampaignID      uint          `json:"campaign_id"`
	FoodType        string        `json:"food_type"`
	Description     string        `json:"description"`
	Quantity        int           `json:"quantity"`
	Weight          int           `json:"weight"`
	ExpiredDate     string        `json:"expired_date"`
	PickUp          int           `json:"pickup"`
	AdditionalChips pq.Int64Array `gorm:"type:integer[]" json:"additional_chips"`
	ChipAcquisition int           `json:"chip_acquisition"`
	IsDone          bool          `gorm:"default:false;not null" json:"is_done"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type UserPersonalDonationInput struct {
	FoodType    string `gorm:"binding:required" json:"food_type"`
	Description string `gorm:"binding:required" json:"description"`
	Quantity    int    `gorm:"binding:required" json:"quantity"`
	Weight      int    `gorm:"binding:required" json:"weight"`
	ExpiredDate string `json:"expired_date"`
}

type UserPersonalDonationConfirmationInput struct {
	PickUp          int           `gorm:"binding:required" json:"pickup"`
	AdditionalChips pq.Int64Array `gorm:"type:integer[]" binding:"required" json:"additional_chips"`
}

type CompanyPersonalDonation struct {
	ID              uint          `gorm:"primaryKey" json:"id"`
	CompanyID       uuid.UUID     `json:"company_id"`
	Campaign        Campaign      `gorm:"ForeignKey:CampaignID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CampaignID      uint          `json:"campaign_id"`
	FoodType        string        `json:"food_type"`
	Description     string        `json:"description"`
	Quantity        int           `json:"quantity"`
	Weight          int           `json:"weight"`
	ExpiredDate     string        `json:"expired_date"`
	PickUp          int           `json:"pickup"`
	AdditionalChips pq.Int64Array `gorm:"type:integer[]" json:"additional_chips"`
	ChipAcquisition int           `json:"chip_acquisition"`
	IsDone          bool          `gorm:"default:false;not null" json:"is_done"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type CompanyPersonalDonationInput struct {
	FoodType    string `gorm:"binding:required" json:"food_type"`
	Description string `gorm:"binding:required" json:"description"`
	Quantity    int    `gorm:"binding:required" json:"quantity"`
	Weight      int    `gorm:"binding:required" json:"weight"`
	ExpiredDate string `json:"expired_date"`
}

type CompanyPersonalDonationConfirmationInput struct {
	PickUp          int           `gorm:"binding:required" json:"pickup"`
	AdditionalChips pq.Int64Array `gorm:"type:integer[]" binding:"required" json:"additional_chips"`
}
