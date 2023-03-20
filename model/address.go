package model

import (
	"time"
)

type Address struct {
	ID              uint   `gorm:"primaryKey" json:"id"`
	User            User   `gorm:"ForeignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	UserID          uint   `json:"user_id"`
	Name            string `json:"name"`
	Phone           string `json:"phone"`
	Address         string `json:"address"`
	State           string `json:"state"`
	City            string `json:"city"`
	Disctrict       string `json:"disctrict"`
	ZipCode         string `json:"zip_code"`
	DetailedAddress string `json:"detailed_address"`
	PrimaryAddress  bool   `json:"primary_address"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type AddressInput struct {
	Name            string `gorm:"binding:required" json:"name"`
	Phone           string `gorm:"binding:required" json:"phone"`
	Address         string `gorm:"binding:required" json:"address"`
	State           string `gorm:"binding:required" json:"state"`
	City            string `gorm:"binding:required" json:"city"`
	Disctrict       string `gorm:"binding:required" json:"disctrict"`
	ZipCode         string `gorm:"binding:required" json:"zip_code"`
	DetailedAddress string `gorm:"binding:required" json:"detailed_address"`
	PrimaryAddress  bool   `gorm:"binding:required" json:"primary_address"`
}
