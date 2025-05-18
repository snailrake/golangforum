package handler

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/websocket"
	"golangforum/internal/client"
	"golangforum/internal/usecase"
	"net/http"
)

type ChatHandler struct {
	UC      usecase.ChatUseCase
	AC      *client.AuthClient
	clients map[*websocket.Conn]struct{}
}

func NewChatHandler(uc usecase.ChatUseCase, ac *client.AuthClient) *ChatHandler {
	return &ChatHandler{
		UC:      uc,
		AC:      ac,
		clients: make(map[*websocket.Conn]struct{}),
	}
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

	user, id, err := h.getUser(r)
	if err != nil {
		switch {
		case errors.Is(err, usecase.ErrMissingToken):
			http.Error(w, "missing token", http.StatusBadRequest)
		case errors.Is(err, usecase.ErrInvalidToken),
			errors.Is(err, usecase.ErrInvalidUserID),
			errors.Is(err, usecase.ErrInvalidUsername):
			http.Error(w, err.Error(), http.StatusUnauthorized)
		default:
			http.Error(w, "internal error", http.StatusInternalServerError)
		}
		conn.Close()
		return
	}

	h.clients[conn] = struct{}{}
	defer func() {
		delete(h.clients, conn)
		conn.Close()
	}()
	h.UC.HandleConnection(conn, user, id, h.clients)
}

func (h *ChatHandler) getUser(r *http.Request) (string, int, error) {
	t := r.URL.Query().Get("token")
	if t == "" {
		return "", 0, usecase.ErrMissingToken
	}
	claims, err := h.AC.VerifyToken(t)
	if err != nil {
		return "", 0, usecase.ErrInvalidToken
	}
	uid, ok := claims["user_id"].(float64)
	if !ok {
		return "", 0, usecase.ErrInvalidUserID
	}
	name, ok := claims["username"].(string)
	if !ok {
		return "", 0, usecase.ErrInvalidUsername
	}
	return name, int(uid), nil
}
