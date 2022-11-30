package routes

import (
	"Stage2Backend/handlers"
	"Stage2Backend/pkg/middleware"
	"Stage2Backend/pkg/mysql"
	"Stage2Backend/repositories"

	"github.com/gorilla/mux"
)

func ToppingRoutes(r *mux.Router) {
	ToppingRepository := repositories.RepositoryTopping(mysql.DB)
	h := handlers.HandlerTopping(ToppingRepository)

	r.HandleFunc("/toppings", h.FindToppings).Methods("GET")
	r.HandleFunc("/topping/{id}", h.GetTopping).Methods("GET")
	r.HandleFunc("/topping", middleware.Auth(middleware.UploadFile(h.CreateTopping))).Methods("POST")
	r.HandleFunc("/topping/{id}", middleware.Auth(middleware.UploadFile(h.UpdateTopping))).Methods("PATCH")
	r.HandleFunc("/topping/{id}", middleware.Auth(h.DeleteTopping)).Methods("DELETE")
}
