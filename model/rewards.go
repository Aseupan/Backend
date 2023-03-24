package model

type Rewards struct {
	ID       uint   `json:"id" gorm:"primaryKey"`
	Name     string `json:"name"`
	Points   int    `json:"points"`
	Quantity int    `json:"quantity"`
}

type RewardsInput struct {
	RewardID int `json:"reward_id"`
}
