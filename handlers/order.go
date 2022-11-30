package handlers

import (
	dto "Stage2Backend/dto/result"
	transactiondto "Stage2Backend/dto/transaction"
	"Stage2Backend/models"
	"Stage2Backend/repositories"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/mux"
)

type handlerOrder struct {
	OrderRepository repositories.OrderRepository
}

func HandlerOrder(OrderRepository repositories.OrderRepository) *handlerOrder {
	return &handlerOrder{OrderRepository}
}

func (h *handlerOrder) FindOrders(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	orders, err := h.OrderRepository.FindOrders()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: http.StatusOK, Data: orders}
	json.NewEncoder(w).Encode(response)
}

func (h *handlerOrder) GetOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	var order models.Order
	order, err := h.OrderRepository.GetOrder(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: http.StatusOK, Data: order}
	json.NewEncoder(w).Encode(response)
}

func (h *handlerOrder) CreateOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	/*
		orderInfo := r.Context().Value("orderInfo").(jwt.MapClaims)
		order_ID := int(orderInfo["id"].(float64))

		fmt.Println(order_ID)*/

	request := new(transactiondto.OrderRequest)
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Status: "failed", Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	validation := validator.New()
	err := validation.Struct(request)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Code: http.StatusInternalServerError, Status: "failed1", Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	topping, _ := h.OrderRepository.FindToppingId(request.ToppingID)

	order := models.Order{
		ProductID: request.ProductID,
		Qty:       request.Qty,
		Topping:   topping,
	}
	//fmt.Println(h.OrderRepository)
	order, err = h.OrderRepository.CreateOrder(order)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Code: http.StatusInternalServerError, Status: "failed2", Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	order, _ = h.OrderRepository.GetOrder(order.ID)

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: http.StatusOK, Status: "success", Data: order}
	json.NewEncoder(w).Encode(response)

}

func (h *handlerOrder) UpdateOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	/*
		productID, _ := strconv.Atoi(r.FormValue("product"))
		qty, _ := strconv.Atoi(r.FormValue("qty"))
		//topps, _ := strconv.Atoi([]r.FormValue("product"))

		request := transactiondto.OrderRequest{
			ProductID:     productID,
			Qty:    qty,
			ToppingID: r.FormValue("toppings"),
		}*/

	request := new(transactiondto.OrderRequest)

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Status: "failed", Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	validation := validator.New()
	err := validation.Struct(request)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Code: http.StatusInternalServerError, Status: "failed1", Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	order, err := h.OrderRepository.GetOrder(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Status: "failed", Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	fmt.Println(request)
	if request.ProductID != 0 {
		order.ProductID = request.ProductID
	}

	if request.Qty != 0 {
		order.Qty = request.Qty
	}

	data, err := h.OrderRepository.UpdateOrder(order)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Code: http.StatusInternalServerError, Status: "failed", Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: http.StatusOK, Status: "success", Data: data}
	json.NewEncoder(w).Encode(response)

}

func (h *handlerOrder) DeleteOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userInfo := r.Context().Value("userInfo").(jwt.MapClaims)
	userRole := userInfo["role"]

	//fmt.Println(userInfo)
	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	if userRole == "admin" {

		order, err := h.OrderRepository.GetOrder(id)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			response := dto.ErrorResult{Code: http.StatusBadRequest, Status: "failed", Message: err.Error()}
			json.NewEncoder(w).Encode(response)
			return
		}

		data, err := h.OrderRepository.DeleteOrder(order)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			response := dto.ErrorResult{Code: http.StatusInternalServerError, Status: "failed1", Message: err.Error()}
			json.NewEncoder(w).Encode(response)
			return
		}

		w.WriteHeader(http.StatusOK)
		response := dto.SuccessResult{Code: http.StatusOK, Status: "success", Data: data.ID}
		json.NewEncoder(w).Encode(response)

	} else {

		w.WriteHeader(http.StatusUnauthorized)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Status: "failed", Message: "unauthorized"}
		json.NewEncoder(w).Encode(response)
		return
	}
}
