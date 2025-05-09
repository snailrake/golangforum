package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/websocket"
	"golangforum/internal/client"
	"golangforum/internal/usecase"
)

type ChatHandler struct {
	UC *usecase.ChatUseCase
	AC *client.AuthClient
}

func NewChatHandler(uc *usecase.ChatUseCase, ac *client.AuthClient) *ChatHandler {
	return &ChatHandler{UC: uc, AC: ac}
}

var upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}

// GetAllMessages godoc
// @Summary Получить все сообщения чата
// @Description Возвращает список всех сообщений чата
// @Tags Чат
// @Accept json
// @Produce json
// @Success 200 {array} string "OK"
// @Failure 500 {object} string "Ошибка сервера"
// @Router /chat/messages [get]
func (h *ChatHandler) GetAllMessages(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	msgs, err := h.UC.GetAllMessages()
	if err != nil {
		http.Error(w, "could not fetch messages", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(msgs)
}

// ServeWS godoc
// @Summary Установить WebSocket соединение
// @Description Устанавливает WebSocket соединение для чата
// @Tags Чат
// @Accept json
// @Produce json
// @Param token query string true "Токен для аутентификации"
// @Success 101 {object} string "WebSocket соединение установлено"
// @Failure 400 {object} string "Неверный запрос"
// @Failure 500 {object} string "Ошибка сервера"
// @Router /chat [get]
func (h *ChatHandler) ServeWS(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "upgrade failed", http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	user, id, err := h.getUser(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	h.UC.HandleConnection(conn, user, id)
}

func (h *ChatHandler) getUser(r *http.Request) (string, int, error) {
	t := r.URL.Query().Get("token")
	if t == "" {
		return "", 0, errors.New("token required")
	}
	claims, err := h.AC.VerifyToken(t)
	if err != nil {
		return "", 0, err
	}
	uid, ok := claims["user_id"].(float64)
	if !ok {
		return "", 0, errors.New("invalid user_id")
	}
	name, ok := claims["username"].(string)
	if !ok {
		return "", 0, errors.New("invalid username")
	}
	return name, int(uid), nil
}
