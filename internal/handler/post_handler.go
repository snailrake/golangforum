package handler

import (
	"encoding/json"
	"fmt"
	"golangforum/internal/domain"
	"golangforum/internal/usecase"
	"net/http"
	"strconv"
)

type PostHandler struct {
	UseCase *usecase.PostUseCase
}

func NewPostHandler(uc *usecase.PostUseCase) *PostHandler {
	return &PostHandler{UseCase: uc}
}

func (h *PostHandler) Create(w http.ResponseWriter, r *http.Request) {
	var p domain.Post
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	if err := h.UseCase.Create(&p); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "post created"})
}

func (h *PostHandler) GetByTopic(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("topic_id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		errorMessage := fmt.Sprintf("invalid topic_id: %v", err) // Добавляем текст ошибки
		http.Error(w, errorMessage, http.StatusBadRequest)
		return
	}
	posts, err := h.UseCase.GetByTopic(id)
	if err != nil {
		http.Error(w, "could not fetch posts", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}

func (h *PostHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	posts, err := h.UseCase.GetAll()
	if err != nil {
		http.Error(w, "could not fetch posts", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}
