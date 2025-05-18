package handler

import (
	"encoding/json"
	"errors"
	"golangforum/internal/model"
	"net/http"
	"strconv"

	"golangforum/internal/usecase"
)

type TopicHandler struct {
	UseCase usecase.TopicUseCase
}

func NewTopicHandler(uc usecase.TopicUseCase) *TopicHandler {
	return &TopicHandler{UseCase: uc}
}

// Create godoc
// @Summary Создать тему
// @Description Создает новую тему с заголовком и описанием, используя структуру Topic
// @Tags Темы
// @Accept json
// @Produce json
// @Param topic body model.Topic true "Данные темы"
// @Success 201 {object} map[string]string "Тема успешно создана"
// @Failure 400 {object} map[string]string "Неверный запрос"
// @Failure 500 {object} map[string]string "Ошибка сервера"
// @Router /topics/create [post]
func (h *TopicHandler) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "use POST", http.StatusMethodNotAllowed)
		return
	}
	var topic model.Topic
	if err := json.NewDecoder(r.Body).Decode(&topic); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if err := topic.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.UseCase.Create(topic.Title, topic.Description); err != nil {
		switch {
		case errors.Is(err, usecase.ErrInvalidTopicData):
			http.Error(w, "invalid topic data", http.StatusBadRequest)
		default:
			http.Error(w, "internal error", http.StatusInternalServerError)
		}
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "topic created"})
}

// GetAll godoc
// @Summary Получить все темы
// @Description Возвращает список всех тем форума
// @Tags Темы
// @Accept json
// @Produce json
// @Success 200 {array} model.Topic "Список тем"
// @Failure 500 {object} map[string]string "Ошибка сервера"
// @Router /topics [get]
func (h *TopicHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "use GET", http.StatusMethodNotAllowed)
		return
	}
	topics, err := h.UseCase.GetAll()
	if err != nil {
		http.Error(w, "could not fetch topics", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(topics)
}

// Delete godoc
// @Summary Удалить тему
// @Description Удаляет тему по ее ID
// @Tags Темы
// @Accept json
// @Produce json
// @Param id query int true "ID темы"
// @Success 200 {object} map[string]string "Тема успешно удалена"
// @Failure 400 {object} map[string]string "Неверный id"
// @Failure 500 {object} map[string]string "Ошибка сервера"
// @Router /topics/delete [delete]
func (h *TopicHandler) Delete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "use DELETE", http.StatusMethodNotAllowed)
		return
	}
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	if err := h.UseCase.Delete(id); err != nil {
		switch {
		case errors.Is(err, usecase.ErrTopicNotFound):
			http.Error(w, "topic not found", http.StatusBadRequest)
		default:
			http.Error(w, "internal error", http.StatusInternalServerError)
		}
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"topic deleted"}`))
}
