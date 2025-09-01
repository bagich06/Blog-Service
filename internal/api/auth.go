package api

import (
	"blog/internal/auth"
	"blog/internal/models"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

func (api *api) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var req models.RegisterRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user := models.User{
		Name:     req.Name,
		Phone:    req.Phone,
		Password: req.Password,
	}

	userID, err := api.db.CreateUser(user)
	if err != nil {
		http.Error(w, "User already exists or database error", http.StatusConflict)
		return
	}

	token, err := auth.GenerateToken(userID, req.Phone)
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	user.ID = userID
	user.Password = ""

	response := models.AuthResponse{
		Token: token,
		User:  user,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Пользователь успешно зарегистрирован",
		"user_id": userID,
	})
}

func (api *api) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user, err := api.db.GetUserByPhone(req.Phone)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	ctx := r.Context()
	err = api.loginLimiter.CheckLoginLimit(ctx, user.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusTooManyRequests)
		return
	}

	if user.Password != req.Password {
		err = api.loginLimiter.RecordFailedLogin(ctx, user.ID)
		if err != nil {
			log.Printf("Failed to record failed login: %v", err)
		}

		http.Error(w, "Incorrect Password", http.StatusUnauthorized)
		return
	}

	err = api.loginLimiter.RecordSuccessfulLogin(ctx, user.ID)
	if err != nil {
		log.Printf("Failed to reset login attempts: %v", err)
	}

	token, err := auth.GenerateToken(user.ID, user.Phone)
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	user.Password = ""

	response := models.AuthResponse{
		Token: token,
		User:  user,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

	json.NewEncoder(w).Encode(map[string]string{
		"message": "Вы успешно вошли в систему",
		"user_id": strconv.Itoa(user.ID),
	})
}
