package models

type Transaction struct {
	ID     int                  `json:"id" gorm:"primary_key:auto_increment"`
	UserID int                  `json:"user_id"`
	User   UsersProfileResponse `json:"user"`
	Status string               `json:"status"  gorm:"type:varchar(255)"`
	Order  []OrderResponse      `json:"order" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type TransactionResponse struct {
	ID     int    `json:"id"`
	Status string `json:"status"`
}

func (TransactionResponse) TableName() string {
	return "transactions"
}
