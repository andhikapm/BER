package repositories

import (
	"Stage2Backend/models"

	"gorm.io/gorm"
)

type OrderRepository interface {
	FindOrders() ([]models.Order, error)
	GetOrder(ID int) (models.Order, error)
	CreateOrder(order models.Order) (models.Order, error)
	UpdateOrder(order models.Order) (models.Order, error)
	DeleteOrder(order models.Order) (models.Order, error)
	FindToppingId(ToppingID []int) ([]models.ToppingResponse, error)
}

func RepositoryOrder(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) FindOrders() ([]models.Order, error) {
	var orders []models.Order
	err := r.db.Preload("Product").Preload("Topping").Find(&orders).Error

	return orders, err
}

func (r *repository) GetOrder(ID int) (models.Order, error) {
	var order models.Order

	err := r.db.Preload("Product").Preload("Topping").First(&order, ID).Error

	return order, err
}

func (r *repository) CreateOrder(order models.Order) (models.Order, error) {
	err := r.db.Preload("Product").Preload("Topping").Create(&order).Error

	return order, err
}

func (r *repository) UpdateOrder(order models.Order) (models.Order, error) {
	err := r.db.Save(&order).Error

	return order, err
}

func (r *repository) DeleteOrder(order models.Order) (models.Order, error) {
	r.db.Model(&order).Association("Topping").Clear()
	err := r.db.Delete(&order).Error

	return order, err
}

func (r *repository) FindToppingId(ToppingID []int) ([]models.ToppingResponse, error) {
	var toppings []models.ToppingResponse
	err := r.db.Find(&toppings, ToppingID).Error

	return toppings, err
}
