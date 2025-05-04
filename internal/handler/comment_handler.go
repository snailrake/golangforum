package handler

import (
	"encoding/json"
	"errors"
	"golangforum/internal/domain"
	"golangforum/internal/usecase"
	"golangforum/internal/utils"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type CommentHandler struct {
	UseCase *usecase.CommentUseCase
}

func NewCommentHandler(uc *usecase.CommentUseCase) *CommentHandler {
	return &CommentHandler{UseCase: uc}
}

func (h *CommentHandler) Create(w http.ResponseWriter, r *http.Request) {
	var c domain.Comment

	// Извлекаем username из JWT
	username, err := getUsernameFromJWT(r)
	if err != nil {
		http.Error(w, "unauthorized: "+err.Error(), http.StatusUnauthorized)
		return
	}

	// Декодируем тело запроса
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Присваиваем извлеченный username
	c.Username = username
	c.Timestamp = time.Now()

	// Создаем комментарий через use case
	if err := h.UseCase.Create(&c); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "comment created"})
}

func (h *CommentHandler) GetByPost(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("post_id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid post_id", http.StatusBadRequest)
		return
	}
	comments, err := h.UseCase.GetByPost(id)
	if err != nil {
		http.Error(w, "could not fetch comments", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(comments)
}

func getUsernameFromJWT(r *http.Request) (string, error) {
	// Извлекаем JWT из заголовка Authorization
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		return "", errors.New("token is missing")
	}

	// Убираем "Bearer " из токена
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	// Проверяем и парсим токен через VerifyToken
	claims, err := utils.VerifyToken(tokenString)
	if err != nil {
		return "", err
	}

	// Извлекаем username из claims
	username, ok := claims["username"].(string)
	if !ok {
		return "", errors.New("username not found in token")
	}
	return username, nil
}
