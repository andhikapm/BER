package transactiondto

type OrderRequest struct {
	ProductID int   `json:"product" form:"product" gorm:"type: int"`
	Qty       int   `json:"qty" form:"qty" gorm:"type: int"`
	ToppingID []int `json:"toppings" form:"toppings" gorm:"type:int[]"`
}
