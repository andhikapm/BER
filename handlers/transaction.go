package handlers

import (
	dto "Stage2Backend/dto/result"
	transactiondto "Stage2Backend/dto/transaction"
	"fmt"
	"strconv"

	"Stage2Backend/models"
	"Stage2Backend/repositories"
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/mux"
)

type Data struct {
	TransID int
	OrderID []int
}

type handlerTransaction struct {
	TransactionRepository repositories.TransactionRepository
}

func HandlerTransaction(TransactionRepository repositories.TransactionRepository) *handlerTransaction {
	return &handlerTransaction{TransactionRepository}
}

func (h *handlerTransaction) FindTransactions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	transaction, err := h.TransactionRepository.FindTransactions()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Status: "failed", Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	var data []models.Transaction
	for _, sT := range transaction {

		var orderData []models.OrderResponse

		for _, s := range sT.Order {

			orderLoop, _ := h.TransactionRepository.FindTransOrders(s.ID)

			orderRes := models.OrderResponse{
				ID:             orderLoop.ID,
				Transaction_ID: sT.ID,
				ProductID:      orderLoop.ID,
				Product:        orderLoop.Product,
				Qty:            orderLoop.Qty,
				Topping:        orderLoop.Topping,
			}
			orderData = append(orderData, orderRes)
		}

		dataGet := models.Transaction{
			ID:     sT.ID,
			UserID: sT.UserID,
			User:   sT.User,
			Status: sT.Status,
			Order:  orderData,
		}
		data = append(data, dataGet)
	}

	//fmt.Println(data)
	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: http.StatusOK, Status: "success", Data: data}
	json.NewEncoder(w).Encode(response)
}

func (h *handlerTransaction) GetTransaction(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	var transaction models.Transaction
	transaction, err := h.TransactionRepository.GetTransaction(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Status: "failed", Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	var orderData []models.OrderResponse

	for _, s := range transaction.Order {

		orderLoop, _ := h.TransactionRepository.FindTransOrders(s.ID)

		orderRes := models.OrderResponse{
			ID:             orderLoop.ID,
			Transaction_ID: transaction.ID,
			ProductID:      orderLoop.ID,
			Product:        orderLoop.Product,
			Qty:            orderLoop.Qty,
			Topping:        orderLoop.Topping,
		}
		orderData = append(orderData, orderRes)
	}

	data := models.Transaction{
		ID:     transaction.ID,
		UserID: transaction.UserID,
		User:   transaction.User,
		Status: transaction.Status,
		Order:  orderData,
	}

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: http.StatusOK, Status: "success", Data: data}
	json.NewEncoder(w).Encode(response)
}

func (h *handlerTransaction) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userInfo := r.Context().Value("userInfo").(jwt.MapClaims)
	user_ID := int(userInfo["id"].(float64))
	//fmt.Println(user_ID)
	request := new(transactiondto.TransactionRequest)
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

	//fmt.Println(request)

	transaction := models.Transaction{
		UserID: user_ID,
		Status: "pending",
	}

	transaction, err = h.TransactionRepository.CreateTransaction(transaction)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Code: http.StatusInternalServerError, Status: "failed2", Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	transaction, _ = h.TransactionRepository.GetTransaction(transaction.ID)

	var TransOrder []models.Order
	for _, s := range request.Order {

		topping, _ := h.TransactionRepository.FindTransToppingId(s.ToppingID)
		order := models.Order{
			Transaction_ID: transaction.ID,
			ProductID:      s.ProductID,
			Qty:            s.Qty,
			Topping:        topping,
		}
		TransOrder = append(TransOrder, order)
	}

	Ordering, err := h.TransactionRepository.CreateTransOrder(TransOrder)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Code: http.StatusInternalServerError, Status: "failed3", Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	var orderData []models.OrderResponse

	for _, s := range Ordering {

		orderLoop, _ := h.TransactionRepository.FindTransOrders(s.ID)

		orderRes := models.OrderResponse{
			ID:             orderLoop.ID,
			Transaction_ID: transaction.ID,
			ProductID:      orderLoop.ID,
			Product:        orderLoop.Product,
			Qty:            orderLoop.Qty,
			Topping:        orderLoop.Topping,
		}
		orderData = append(orderData, orderRes)
	}

	data := models.Transaction{
		ID:     transaction.ID,
		UserID: transaction.UserID,
		User:   transaction.User,
		Status: transaction.Status,
		Order:  orderData,
	}
	fmt.Println(data)
	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: http.StatusOK, Status: "otw", Data: data}
	json.NewEncoder(w).Encode(response)
}

func (h *handlerTransaction) UpdateTransaction(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	//userInfo := r.Context().Value("userInfo").(jwt.MapClaims)
	//userRole := userInfo["role"]
	//userID := int(userInfo["id"].(float64))

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	request := new(transactiondto.TransactionUpdate)
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

	transaction, err := h.TransactionRepository.GetTransaction(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Status: "failed", Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	if request.Status != "" {
		transaction.Status = request.Status
	}

	transaction, err = h.TransactionRepository.UpdateTransaction(transaction)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Code: http.StatusInternalServerError, Status: "failed", Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	var orderData []models.OrderResponse

	for _, s := range transaction.Order {

		orderLoop, _ := h.TransactionRepository.FindTransOrders(s.ID)

		orderRes := models.OrderResponse{
			ID:             orderLoop.ID,
			Transaction_ID: transaction.ID,
			ProductID:      orderLoop.ID,
			Product:        orderLoop.Product,
			Qty:            orderLoop.Qty,
			Topping:        orderLoop.Topping,
		}
		orderData = append(orderData, orderRes)
	}

	data := models.Transaction{
		ID:     transaction.ID,
		UserID: transaction.UserID,
		User:   transaction.User,
		Status: transaction.Status,
		Order:  orderData,
	}

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: http.StatusOK, Status: "success", Data: data}
	json.NewEncoder(w).Encode(response)

}

func (h *handlerTransaction) DeleteTransaction(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	//userInfo := r.Context().Value("userInfo").(jwt.MapClaims)
	//userRole := userInfo["role"]
	//userID := int(userInfo["id"].(float64))

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	transaction, err := h.TransactionRepository.GetTransaction(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Status: "failed", Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	order, err := h.TransactionRepository.WhereTransOrder(transaction.ID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Status: "failed", Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	order, err = h.TransactionRepository.DeleteTransOrder(order)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Status: "failed", Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	var getID []int

	for _, s := range order {
		getID = append(getID, s.ID)
	}

	transaction, err = h.TransactionRepository.DeleteTransaction(transaction)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Code: http.StatusInternalServerError, Status: "failed", Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	data := Data{
		TransID: transaction.ID,
		OrderID: getID,
	}
	//h.TransactionRepository.DeleteTransOrder(transaction.Order)

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: http.StatusOK, Status: "success", Data: data}
	json.NewEncoder(w).Encode(response)

}

func (h *handlerTransaction) GetMyTrans(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userInfo := r.Context().Value("userInfo").(jwt.MapClaims)
	userID := int(userInfo["id"].(float64))

	//fmt.Println(userID)
	transaction, err := h.TransactionRepository.GetMyTransaction(userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Status: "failed", Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	var data []models.Transaction
	for _, sT := range transaction {

		var orderData []models.OrderResponse

		for _, s := range sT.Order {

			orderLoop, _ := h.TransactionRepository.FindTransOrders(s.ID)

			orderRes := models.OrderResponse{
				ID:             orderLoop.ID,
				Transaction_ID: sT.ID,
				ProductID:      orderLoop.ID,
				Product:        orderLoop.Product,
				Qty:            orderLoop.Qty,
				Topping:        orderLoop.Topping,
			}
			orderData = append(orderData, orderRes)
		}

		dataGet := models.Transaction{
			ID:     sT.ID,
			UserID: sT.UserID,
			User:   sT.User,
			Status: sT.Status,
			Order:  orderData,
		}
		data = append(data, dataGet)
	}

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: http.StatusOK, Status: "success", Data: data}
	json.NewEncoder(w).Encode(response)
}
