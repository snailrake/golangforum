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

type CommentHandler struct {
	uc   usecase.CommentUseCase
	auth *client.AuthClient
}

func NewCommentHandler(uc usecase.CommentUseCase, auth *client.AuthClient) *CommentHandler {
	return &CommentHandler{uc: uc, auth: auth}
}

// Create godoc
// @Summary Создать новый комментарий
// @Description Создает новый комментарий для заданного поста
// @Tags Комментарии
// @Accept json
// @Produce json
// @Param comment body model.Comment true "Комментарий"
// @Success 201 {object} map[string]string "Комментарий создан"
// @Failure 400 {object} map[string]string "Неверный запрос"
// @Failure 401 {object} map[string]string "Не авторизован"
// @Failure 500 {object} map[string]string "Ошибка сервера"
// @Router /comments/create [post]
func (h *CommentHandler) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	user, err := h.auth.GetUsername(r)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	var c model.Comment
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if err := h.uc.Create(user, &c); err != nil {
		switch {
		case errors.Is(err, usecase.ErrInvalidCommentData):
			http.Error(w, "invalid comment", http.StatusBadRequest)
		default:
			http.Error(w, "internal error", http.StatusInternalServerError)
		}
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "created"})
}

// GetByPost godoc
// @Summary Получить все комментарии для поста
// @Description Получает список всех комментариев для заданного поста
// @Tags Комментарии
// @Accept json
// @Produce json
// @Param post_id query int true "ID поста"
// @Success 200 {array} model.Comment "Список комментариев"
// @Failure 400 {object} map[string]string "Неверный запрос"
// @Failure 500 {object} map[string]string "Ошибка сервера"
// @Router /comments [get]
func (h *CommentHandler) GetByPost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	id, err := strconv.Atoi(r.URL.Query().Get("post_id"))
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	comments, err := h.uc.GetByPost(id)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(comments)
}

// Delete godoc
// @Summary Удалить комментарий
// @Description Удаляет комментарий по ID
// @Tags Комментарии
// @Accept json
// @Produce json
// @Param comment_id query int true "ID комментария"
// @Success 200 {object} map[string]string "Комментарий удален"
// @Failure 400 {object} map[string]string "Неверный запрос"
// @Failure 500 {object} map[string]string "Ошибка сервера"
// @Router /comments/delete [delete]
func (h *CommentHandler) Delete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	id, err := strconv.Atoi(r.URL.Query().Get("comment_id"))
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	if err := h.uc.Delete(id); err != nil {
		switch {
		case errors.Is(err, usecase.ErrCommentNotFound):
			http.Error(w, "comment not found", http.StatusBadRequest)
		default:
			http.Error(w, "internal error", http.StatusInternalServerError)
		}
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "deleted"})
}
