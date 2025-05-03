package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"golangforum/internal/domain"
	"golangforum/internal/usecase"
	"golangforum/internal/utils"
)

type AuthHandler struct {
	UseCase *usecase.AuthUseCase
}

func NewAuthHandler(uc *usecase.AuthUseCase) *AuthHandler {
	return &AuthHandler{UseCase: uc}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Используйте POST", http.StatusMethodNotAllowed)
		return
	}
	var u domain.User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, "Неверный запрос", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if err := h.UseCase.Register(&u); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]any{
		"message": "Пользователь зарегистрирован",
		"user_id": u.ID,
	})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Используйте POST", http.StatusMethodNotAllowed)
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

	access, refresh, err := h.UseCase.Login(req.Username, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{
		"access_token":  access,
		"refresh_token": refresh,
	})
}

func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Используйте POST", http.StatusMethodNotAllowed)
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
	auth := r.Header.Get("Authorization")
	if auth == "" {
		http.Error(w, "нет токена", http.StatusUnauthorized)
		return
	}
	parts := strings.Split(auth, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		http.Error(w, "неверный формат токена", http.StatusUnauthorized)
		return
	}

	claims, err := utils.VerifyToken(parts[1])
	if err != nil {
		http.Error(w, "неверный или просроченный токен", http.StatusUnauthorized)
		return
	}
	username, ok := claims["username"].(string)
	if !ok {
		http.Error(w, "ошибка токена", http.StatusUnauthorized)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{
		"message": fmt.Sprintf("Добро пожаловать, %s!", username),
	})
}
