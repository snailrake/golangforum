package handler

import (
	"encoding/json"
	"fmt"
	"golangforum/internal/domain"
	"golangforum/internal/usecase"
	"golangforum/internal/utils"
	"net/http"
	"strings"
)

type AuthHandler struct {
	UseCase *usecase.AuthUseCase
}

func NewAuthHandler(uc *usecase.AuthUseCase) *AuthHandler {
	return &AuthHandler{UseCase: uc}
}

func (h *AuthHandler) Register(response http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		http.Error(response, "Используйте метод POST", http.StatusMethodNotAllowed)
		return
	}

	var user domain.User
	if err := json.NewDecoder(request.Body).Decode(&user); err != nil {
		http.Error(response, "Неверный запрос", http.StatusBadRequest)
		return
	}
	defer request.Body.Close()

	if err := h.UseCase.Register(&user); err != nil {
		http.Error(response, err.Error(), http.StatusBadRequest)
		return
	}
	response.WriteHeader(http.StatusCreated)
	json.NewEncoder(response).Encode(map[string]any{
		"message": "Пользователь зарегистрирован",
		"user_id": user.ID,
	})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Используйте метод POST", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Неверный запрос", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	accessToken, refreshToken, err := h.UseCase.Login(req.Username, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Используйте метод POST", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Неверный запрос", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	newAccess, newRefresh, err := h.UseCase.RefreshToken(req.RefreshToken)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{
		"access_token":  newAccess,
		"refresh_token": newRefresh,
	})
}

func (h *AuthHandler) Protected(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Нет токена", http.StatusUnauthorized)
		return
	}
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		http.Error(w, "Неверный формат токена", http.StatusUnauthorized)
		return
	}
	tokenString := parts[1]

	claims, err := utils.VerifyToken(tokenString)
	if err != nil {
		http.Error(w, "Неверный или просроченный токен", http.StatusUnauthorized)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{
		"message": fmt.Sprintf("Добро пожаловать, %s!", claims.Username),
	})
}
