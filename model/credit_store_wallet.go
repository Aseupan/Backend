package model

import "github.com/google/uuid"

type CreditStoreWallet struct {
	User   User      `gorm:"ForeignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	UserID uuid.UUID `json:"user_id"`
	Points int       `json:"points"`
	Price  int       `json:"price"`
}

type CreditStoreWalletInput struct {
	ID int `json:"id"`
}

type CartItem struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	Quantity int       `json:"quantity"`
	Price    int       `json:"price"`
}

type CreditCartItem struct {
	CreditStoreWallet
	Quantity int `json:"quantity"`
}

// type ViewCart struct {
// 	Points     int `json:"points"`
// 	Price      int `json:"price"`
// 	Quantity   int `json:"quantity"`
// 	TotalPrice int `json:"total_price"`
// }

// cart := []CreditCartItem{
// 	{
// 		CreditStoreWallet: CreditStoreWallet{
// 			UserID: user1.ID,
// 			Points: 100,
// 			Price:  50,
// 		},
// 		Quantity: 2,
// 	},
// 	{
// 		CreditStoreWallet: CreditStoreWallet{
// 			UserID: user2.ID,
// 			Points: 50,
// 			Price:  25,
// 		},
// 		Quantity: 1,
// 	},
// }

// 	var totalPrice int
//  for _, item := range cart {
//     totalPrice += item.Price * item.Quantity
//  }
