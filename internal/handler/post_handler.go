package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"golangforum/internal/client"
	"golangforum/internal/model"
	"golangforum/internal/usecase"
)

type PostHandler struct {
	UseCase    usecase.PostUseCase
	AuthClient *client.AuthClient
}

func NewPostHandler(uc usecase.PostUseCase, authClient *client.AuthClient) *PostHandler {
	return &PostHandler{UseCase: uc, AuthClient: authClient}
}

// Create godoc
// @Summary Создать новый пост
// @Description Создает новый пост в указанной теме
// @Tags Посты
// @Accept json
// @Produce json
// @Param post body model.Post true "Данные поста"
// @Success 201 {object} map[string]string "Пост успешно создан"
// @Failure 400 {object} map[string]string "Неверный запрос"
// @Failure 401 {object} map[string]string "Не авторизован"
// @Failure 500 {object} map[string]string "Ошибка сервера"
// @Router /posts/create [post]
func (h *PostHandler) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "use POST", http.StatusMethodNotAllowed)
		return
	}
	var p model.Post
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	username, err := h.AuthClient.GetUsername(r)
	if err != nil {
		http.Error(w, "unauthorized: "+err.Error(), http.StatusUnauthorized)
		return
	}
	if err := h.UseCase.Create(username, &p); err != nil {
		switch {
		case errors.Is(err, usecase.ErrInvalidPostData):
			http.Error(w, "invalid post data", http.StatusBadRequest)
		default:
			http.Error(w, "internal error", http.StatusInternalServerError)
		}
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "post created"})
}

// GetByTopic godoc
// @Summary Получить все посты по теме
// @Description Возвращает список всех постов для указанной темы
// @Tags Посты
// @Accept json
// @Produce json
// @Param topic_id query int true "ID темы"
// @Success 200 {array} model.Post "Список постов"
// @Failure 400 {object} map[string]string "Неверный topic_id"
// @Failure 500 {object} map[string]string "Ошибка сервера"
// @Router /posts [get]
func (h *PostHandler) GetByTopic(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "use GET", http.StatusMethodNotAllowed)
		return
	}
	id, err := strconv.Atoi(r.URL.Query().Get("topic_id"))
	if err != nil {
		http.Error(w, "invalid topic_id", http.StatusBadRequest)
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

// GetAll godoc
// @Summary Получить все посты
// @Description Возвращает список всех постов
// @Tags Посты
// @Accept json
// @Produce json
// @Success 200 {array} model.Post "Список всех постов"
// @Failure 500 {object} map[string]string "Ошибка сервера"
// @Router /posts/all [get]
func (h *PostHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "use GET", http.StatusMethodNotAllowed)
		return
	}
	posts, err := h.UseCase.GetAll()
	if err != nil {
		http.Error(w, "could not fetch posts", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}

// Delete godoc
// @Summary Удалить пост
// @Description Удаляет пост по его ID
// @Tags Посты
// @Accept json
// @Produce json
// @Param post_id query int true "ID поста"
// @Success 200 {object} map[string]string "Пост успешно удален"
// @Failure 400 {object} map[string]string "Неверный post_id"
// @Failure 500 {object} map[string]string "Ошибка сервера"
// @Router /posts/delete [delete]
func (h *PostHandler) Delete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "use DELETE", http.StatusMethodNotAllowed)
		return
	}
	id, err := strconv.Atoi(r.URL.Query().Get("post_id"))
	if err != nil {
		http.Error(w, "invalid post_id", http.StatusBadRequest)
		return
	}
	if err := h.UseCase.Delete(id); err != nil {
		switch {
		case errors.Is(err, usecase.ErrPostNotFound):
			http.Error(w, "post not found", http.StatusBadRequest)
		default:
			http.Error(w, "internal error", http.StatusInternalServerError)
		}
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "post deleted"})
}
