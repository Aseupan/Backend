package model

import "github.com/google/uuid"

type CreditStoreCart struct {
	User          User        `gorm:"ForeignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	UserID        uuid.UUID   `json:"user_id"`
	CreditStore   CreditStore `gorm:"ForeignKey:CreditStoreID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CreditStoreID uint        `json:"credit_store_id"`
	Points        int         `json:"points"`
	Price         int         `json:"price"`
	Quantity      int         `json:"quantity"`
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
