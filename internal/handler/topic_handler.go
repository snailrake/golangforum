package handler

import (
	"encoding/json"
	"golangforum/internal/domain"
	"golangforum/internal/usecase"
	"net/http"
)

type TopicHandler struct {
	UseCase *usecase.TopicUseCase
}

func NewTopicHandler(uc *usecase.TopicUseCase) *TopicHandler {
	return &TopicHandler{UseCase: uc}
}

func (h *TopicHandler) Create(w http.ResponseWriter, r *http.Request) {
	var t domain.Topic
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	if err := h.UseCase.Create(&t); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "topic created"})
}

func (h *TopicHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	topics, err := h.UseCase.GetAll()
	if err != nil {
		http.Error(w, "could not fetch topics", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(topics)
}
