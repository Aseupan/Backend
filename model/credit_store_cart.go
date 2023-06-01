package model

import "github.com/google/uuid"

type CreditStoreCart struct {
	// jika tidak dipakai, maka hapus User User, agar tidak menggunakan foreign key
	User          User      `gorm:"ForeignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	UserID        uuid.UUID `gorm:"null" json:"user_id"`
	CompanyID     uuid.UUID `gorm:"null" json:"company_id"`
	CreditStoreID uint      `json:"credit_store_id"`
	Points        int       `json:"points"`
	Price         int       `json:"price"`
	Quantity      int       `json:"quantity"`
}

type CreditStoreCartInput struct {
	ID int `json:"id"`
}

type CreditStoreUpdateQuantityInput struct {
	Quantity int `json:"quantity"`
}

type CreditStorePaymentInput struct {
	PaymentMethod int `json:"payment_method"`
}
