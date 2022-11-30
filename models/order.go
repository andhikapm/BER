package models

type Order struct {
	ID             int                 `json:"id" gorm:"primary_key:auto_increment"`
	Transaction_ID int                 `json:"transaction_id" form:"transaction_id"`
	Transaction    TransactionResponse `json:"transaction" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	ProductID      int                 `json:"product_id"`
	Product        ProductResponse     `json:"product"`
	Qty            int                 `json:"qty" form:"qty"`
	Topping        []ToppingResponse   `json:"topping" gorm:"many2many:order_topping;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type OrderResponse struct {
	ID             int               `json:"id"`
	Transaction_ID int               `json:"-"`
	ProductID      int               `json:"product_id"`
	Product        ProductResponse   `json:"product"`
	Qty            int               `json:"qty"`
	Topping        []ToppingResponse `json:"topping" gorm:"many2many:order_topping"`
}

func (OrderResponse) TableName() string {
	return "orders"
}
