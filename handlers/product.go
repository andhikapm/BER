package handlers

import (
	productdto "Stage2Backend/dto/product"
	dto "Stage2Backend/dto/result"
	"Stage2Backend/models"
	"Stage2Backend/repositories"
	"encoding/json"
	"net/http"
	"os"
	"strconv"

	"github.com/go-playground/validator/v10"

	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/mux"
)

//var path_file = os.Getenv("PATH_FILE")

type handlerProduct struct {
	ProductRepository repositories.ProductRepository
}

func HandlerProduct(ProductRepository repositories.ProductRepository) *handlerProduct {
	return &handlerProduct{ProductRepository}
}

func (h *handlerProduct) FindProducts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	products, err := h.ProductRepository.FindProducts()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	for i, p := range products {
		products[i].Image = os.Getenv("PATH_FILE") + p.Image
	}

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: http.StatusOK, Data: products}
	json.NewEncoder(w).Encode(response)
}

func (h *handlerProduct) GetProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	var product models.Product
	product, err := h.ProductRepository.GetProduct(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	product.Image = os.Getenv("PATH_FILE") + product.Image

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: http.StatusOK, Data: product}
	json.NewEncoder(w).Encode(response)
}

func (h *handlerProduct) CreateProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userInfo := r.Context().Value("userInfo").(jwt.MapClaims)
	userRole := userInfo["role"]

	if userRole == "admin" {
		dataContex := r.Context().Value("dataFile")
		filename := dataContex.(string)

		price, _ := strconv.Atoi(r.FormValue("price"))
		request := productdto.ProductRequest{
			Title: r.FormValue("title"),
			Price: price,
		}

		validation := validator.New()
		err := validation.Struct(request)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			response := dto.ErrorResult{Code: http.StatusInternalServerError, Status: "failed", Message: err.Error()}
			json.NewEncoder(w).Encode(response)
			return
		}

		product := models.Product{
			Title: request.Title,
			Price: request.Price,
			Image: filename,
		}

		product, err = h.ProductRepository.CreateProduct(product)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			response := dto.ErrorResult{Code: http.StatusInternalServerError, Status: "failed", Message: err.Error()}
			json.NewEncoder(w).Encode(response)
			return
		}

		product, _ = h.ProductRepository.GetProduct(product.ID)

		w.WriteHeader(http.StatusOK)
		response := dto.SuccessResult{Code: http.StatusOK, Status: "success", Data: product}
		json.NewEncoder(w).Encode(response)

	} else {

		w.WriteHeader(http.StatusUnauthorized)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Status: "failed", Message: "unauthorized"}
		json.NewEncoder(w).Encode(response)
		return
	}
}

func (h *handlerProduct) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userInfo := r.Context().Value("userInfo").(jwt.MapClaims)
	userRole := userInfo["role"]

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	if userRole == "admin" {

		price, _ := strconv.Atoi(r.FormValue("price"))

		request := productdto.ProductRequest{
			Title: r.FormValue("title"),
			Price: price,
		}

		product, err := h.ProductRepository.GetProduct(int(id))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			response := dto.ErrorResult{Code: http.StatusBadRequest, Status: "failed", Message: err.Error()}
			json.NewEncoder(w).Encode(response)
			return
		}

		dataContex := r.Context().Value("dataFile")
		filename := dataContex.(string)

		if request.Title != "" {
			product.Title = request.Title
		}

		if r.FormValue("price") != "" {
			product.Price = request.Price
		}

		product.Image = filename

		data, err := h.ProductRepository.UpdateProduct(product)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			response := dto.ErrorResult{Code: http.StatusInternalServerError, Status: "failed", Message: err.Error()}
			json.NewEncoder(w).Encode(response)
			return
		}

		w.WriteHeader(http.StatusOK)
		response := dto.SuccessResult{Code: http.StatusOK, Status: "success", Data: data}
		json.NewEncoder(w).Encode(response)

	} else {

		w.WriteHeader(http.StatusUnauthorized)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Status: "failed", Message: "unauthorized"}
		json.NewEncoder(w).Encode(response)
		return
	}
}

func (h *handlerProduct) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userInfo := r.Context().Value("userInfo").(jwt.MapClaims)
	userRole := userInfo["role"]

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	if userRole == "admin" {

		product, err := h.ProductRepository.GetProduct(id)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			response := dto.ErrorResult{Code: http.StatusBadRequest, Status: "failed", Message: err.Error()}
			json.NewEncoder(w).Encode(response)
			return
		}

		data, err := h.ProductRepository.DeleteProduct(product)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			response := dto.ErrorResult{Code: http.StatusInternalServerError, Status: "failed", Message: err.Error()}
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
