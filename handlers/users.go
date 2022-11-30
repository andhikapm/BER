package handlers

import (
	dto "Stage2Backend/dto/result"
	usersdto "Stage2Backend/dto/users"
	"Stage2Backend/models"
	"Stage2Backend/pkg/bcrypt"
	"Stage2Backend/repositories"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/mux"
)

type handlerUser struct {
	UserRepository repositories.UserRepository
}

func HandlerUser(UserRepository repositories.UserRepository) *handlerUser {
	return &handlerUser{UserRepository}
}

func (h *handlerUser) FindUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	users, err := h.UserRepository.FindUsers()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err.Error())
	}

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: http.StatusOK, Status: "success", Data: users}
	json.NewEncoder(w).Encode(response)
}

func (h *handlerUser) UpdateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userInfo := r.Context().Value("userInfo").(jwt.MapClaims)
	userID := int(userInfo["id"].(float64))

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	if userID == id {
		request := usersdto.UpdateUserRequest{
			Name:     r.FormValue("name"),
			Email:    r.FormValue("email"),
			Password: r.FormValue("password"),
		}

		user, err := h.UserRepository.GetUser(int(id))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			response := dto.ErrorResult{Code: http.StatusBadRequest, Status: "failed", Message: err.Error()}
			json.NewEncoder(w).Encode(response)
			return
		}

		dataContex := r.Context().Value("dataFile")
		filename := dataContex.(string)

		if request.Name != "" {
			user.Name = request.Name
		}

		if request.Email != "" {
			user.Email = request.Email
		}

		if request.Password != "" {
			password, err := bcrypt.HashingPassword(request.Password)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				response := dto.ErrorResult{Code: http.StatusInternalServerError, Status: "failed", Message: err.Error()}
				json.NewEncoder(w).Encode(response)
			}
			user.Password = password
		}

		user.Image = filename

		data, err := h.UserRepository.UpdateUser(user)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			response := dto.ErrorResult{Code: http.StatusInternalServerError, Status: "failed", Message: err.Error()}
			json.NewEncoder(w).Encode(response)
			return
		}

		w.WriteHeader(http.StatusOK)
		response := dto.SuccessResult{Code: http.StatusOK, Status: "success", Data: convertResponse(data)}
		json.NewEncoder(w).Encode(response)

	} else {

		w.WriteHeader(http.StatusUnauthorized)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Status: "failed", Message: "unauthorized"}
		json.NewEncoder(w).Encode(response)
		return
	}
}

func (h *handlerUser) DeleteUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userInfo := r.Context().Value("userInfo").(jwt.MapClaims)
	userID := userInfo["id"]
	userRole := userInfo["role"]

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	if (userID == id) || (userRole == "admin") {

		user, err := h.UserRepository.GetUser(id)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			response := dto.ErrorResult{Code: http.StatusBadRequest, Status: "failed", Message: err.Error()}
			json.NewEncoder(w).Encode(response)
			return
		}

		data, err := h.UserRepository.DeleteUser(user)
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

func convertResponse(u models.User) usersdto.UserResponse {
	return usersdto.UserResponse{
		ID:       u.ID,
		Name:     u.Name,
		Email:    u.Email,
		Password: u.Password,
		Image:    u.Image,
	}
}
