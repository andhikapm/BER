package toppingdto

type ToppingRequest struct {
	Title string `json:"name" form:"name" gorm:"type: varchar(255)"`
	Price int    `json:"price" form:"price" gorm:"type: int"`
}
