package routes

import (
	"Stage2Backend/handlers"
	"Stage2Backend/pkg/middleware"
	"Stage2Backend/pkg/mysql"
	"Stage2Backend/repositories"

	"github.com/gorilla/mux"
)

func OrderRoutes(r *mux.Router) {
	OrderRepository := repositories.RepositoryOrder(mysql.DB)
	h := handlers.HandlerOrder(OrderRepository)

	r.HandleFunc("/orders", h.FindOrders).Methods("GET")
	r.HandleFunc("/order/{id}", h.GetOrder).Methods("GET")
	r.HandleFunc("/order", h.CreateOrder).Methods("POST")
	r.HandleFunc("/order/{id}", h.UpdateOrder).Methods("PATCH")
	r.HandleFunc("/order/{id}", middleware.Auth(h.DeleteOrder)).Methods("DELETE")
}
