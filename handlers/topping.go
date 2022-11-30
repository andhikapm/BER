package handlers

import (
	dto "Stage2Backend/dto/result"
	toppingdto "Stage2Backend/dto/topping"
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

type handlerTopping struct {
	ToppingRepository repositories.ToppingRepository
}

func HandlerTopping(ToppingRepository repositories.ToppingRepository) *handlerTopping {
	return &handlerTopping{ToppingRepository}
}

func (h *handlerTopping) FindToppings(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	toppings, err := h.ToppingRepository.FindToppings()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	for i, p := range toppings {
		toppings[i].Image = os.Getenv("PATH_FILE") + p.Image
	}

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: http.StatusOK, Data: toppings}
	json.NewEncoder(w).Encode(response)
}

func (h *handlerTopping) GetTopping(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	topping, err := h.ToppingRepository.GetTopping(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	topping.Image = os.Getenv("PATH_FILE") + topping.Image

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: http.StatusOK, Data: topping}
	json.NewEncoder(w).Encode(response)
}

func (h *handlerTopping) CreateTopping(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userInfo := r.Context().Value("userInfo").(jwt.MapClaims)
	userRole := userInfo["role"]

	if userRole == "admin" {
		dataContex := r.Context().Value("dataFile")
		filename := dataContex.(string)

		price, _ := strconv.Atoi(r.FormValue("price"))
		request := toppingdto.ToppingRequest{
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

		topping := models.Topping{
			Title: request.Title,
			Price: request.Price,
			Image: filename,
		}

		topping, err = h.ToppingRepository.CreateTopping(topping)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			response := dto.ErrorResult{Code: http.StatusInternalServerError, Status: "failed", Message: err.Error()}
			json.NewEncoder(w).Encode(response)
			return
		}

		topping, _ = h.ToppingRepository.GetTopping(topping.ID)

		w.WriteHeader(http.StatusOK)
		response := dto.SuccessResult{Code: http.StatusOK, Status: "success", Data: topping}
		json.NewEncoder(w).Encode(response)
	} else {

		w.WriteHeader(http.StatusUnauthorized)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Status: "failed", Message: "unauthorized"}
		json.NewEncoder(w).Encode(response)
		return
	}
}

func (h *handlerTopping) UpdateTopping(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userInfo := r.Context().Value("userInfo").(jwt.MapClaims)
	userRole := userInfo["role"]

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	if userRole == "admin" {

		price, _ := strconv.Atoi(r.FormValue("price"))

		request := toppingdto.ToppingRequest{
			Title: r.FormValue("title"),
			Price: price,
		}

		topping, err := h.ToppingRepository.GetTopping(int(id))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			response := dto.ErrorResult{Code: http.StatusBadRequest, Status: "failed", Message: err.Error()}
			json.NewEncoder(w).Encode(response)
			return
		}

		dataContex := r.Context().Value("dataFile")
		filename := dataContex.(string)

		if request.Title != "" {
			topping.Title = request.Title
		}

		if r.FormValue("price") != "" {
			topping.Price = request.Price
		}

		topping.Image = filename

		data, err := h.ToppingRepository.UpdateTopping(topping)
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

func (h *handlerTopping) DeleteTopping(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userInfo := r.Context().Value("userInfo").(jwt.MapClaims)
	userRole := userInfo["role"]

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	if userRole == "admin" {

		topping, err := h.ToppingRepository.GetTopping(id)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			response := dto.ErrorResult{Code: http.StatusBadRequest, Status: "failed", Message: err.Error()}
			json.NewEncoder(w).Encode(response)
			return
		}

		data, err := h.ToppingRepository.DeleteTopping(topping)
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
