package repositories

import (
	"Stage2Backend/models"

	"gorm.io/gorm"
)

type TransactionRepository interface {
	FindTransactions() ([]models.Transaction, error)
	GetTransaction(ID int) (models.Transaction, error)
	CreateTransaction(transaction models.Transaction) (models.Transaction, error)
	UpdateTransaction(transaction models.Transaction) (models.Transaction, error)
	DeleteTransaction(transaction models.Transaction) (models.Transaction, error)
	CreateTransOrder(order []models.Order) ([]models.Order, error)
	FindTransToppingId(ToppingID []int) ([]models.ToppingResponse, error)
	FindTransOrders(orderID int) (models.Order, error)
	WhereTransOrder(ID int) ([]models.Order, error)
	DeleteTransOrder(order []models.Order) ([]models.Order, error)
	GetMyTransaction(ID int) ([]models.Transaction, error)
}

func RepositoryTransaction(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) FindTransactions() ([]models.Transaction, error) {
	var transactions []models.Transaction
	err := r.db.Preload("User").Preload("Order").Find(&transactions).Error

	return transactions, err
}

func (r *repository) GetTransaction(ID int) (models.Transaction, error) {
	var transaction models.Transaction
	err := r.db.Preload("User").Preload("Order").First(&transaction, ID).Error

	return transaction, err
}

func (r *repository) CreateTransaction(transaction models.Transaction) (models.Transaction, error) {
	r.db.Model(&transaction).Association("Order").Append(transaction.Order)
	err := r.db.Preload("User").Preload("Order").Create(&transaction).Error

	return transaction, err
}

func (r *repository) UpdateTransaction(transaction models.Transaction) (models.Transaction, error) {
	r.db.Model(&transaction).Association("Order").Replace(transaction.Order)

	err := r.db.Preload("User").Preload("Order").Save(&transaction).Error

	return transaction, err
}

func (r *repository) DeleteTransaction(transaction models.Transaction) (models.Transaction, error) {
	err := r.db.Delete(&transaction).Error

	return transaction, err
}

func (r *repository) CreateTransOrder(order []models.Order) ([]models.Order, error) {
	err := r.db.Preload("Product").Preload("Topping").Create(&order).Error

	return order, err
}

func (r *repository) FindTransToppingId(ToppingID []int) ([]models.ToppingResponse, error) {
	var toppings []models.ToppingResponse
	err := r.db.Find(&toppings, ToppingID).Error

	return toppings, err
}

func (r *repository) FindTransOrders(orderID int) (models.Order, error) {
	var order models.Order
	err := r.db.Preload("Product").Preload("Topping").First(&order, orderID).Error

	return order, err
}

func (r *repository) WhereTransOrder(ID int) ([]models.Order, error) {
	var order []models.Order
	//db.Where("name = ?", "jinzhu").First(&user)
	err := r.db.Preload("Product").Preload("Topping").Where("transaction_id = ?", ID).Find(&order).Error

	return order, err
}

func (r *repository) DeleteTransOrder(order []models.Order) ([]models.Order, error) {
	r.db.Model(&order).Association("Topping").Clear()
	err := r.db.Delete(&order).Error

	return order, err
}

func (r *repository) GetMyTransaction(ID int) ([]models.Transaction, error) {
	var transaction []models.Transaction
	//db.Where("name = ?", "jinzhu").First(&user)
	err := r.db.Preload("User").Preload("Order").Where("user_id = ?", ID).Find(&transaction).Error

	return transaction, err
}
