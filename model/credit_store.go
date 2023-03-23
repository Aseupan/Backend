package model

type CreditStore struct {
	ID     uint `json:"id" gorm:"primaryKey"`
	Points int  `json:"points"`
	Price  int  `json:"price"`
}
