// handler/chat_handler.go
package handler

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"golangforum/internal/domain"
	"golangforum/internal/usecase"
	"golangforum/internal/utils"
	"log"
	"net/http"
	"time"
)

type ChatHandler struct {
	UseCase          *usecase.ChatUseCase
	WebSocketManager *WebSocketManager // Используем уже существующий WebSocketManager
}

func NewChatHandler(uc *usecase.ChatUseCase, wsManager *WebSocketManager) *ChatHandler {
	return &ChatHandler{UseCase: uc, WebSocketManager: wsManager}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Метод для получения всех сообщений
func (h *ChatHandler) GetAllMessages(w http.ResponseWriter, r *http.Request) {
	messages, err := h.UseCase.GetAllMessages() // Получаем все сообщения
	if err != nil {
		http.Error(w, "Error fetching messages", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages) // Отправляем сообщения в ответ
}

func (h *ChatHandler) ServeWS(w http.ResponseWriter, r *http.Request) {
	// Извлекаем токен из query параметров URL
	token := r.URL.Query().Get("token")

	// Для неавторизованных пользователей, устанавливаем имя как 'Гость'
	var username string
	var userIDInt int // Объявляем переменную, которая будет хранить ID пользователя
	if token != "" {
		// Проверяем токен
		claims, err := utils.VerifyToken(token)
		if err != nil {
			log.Println("Error verifying token:", err)
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Извлекаем данные из токена
		userID, ok := claims["user_id"].(float64)
		if !ok {
			log.Println("Invalid token data: user_id not found")
			http.Error(w, "Invalid token data", http.StatusUnauthorized)
			return
		}

		userIDInt = int(userID) // Преобразуем в int
		username, ok = claims["username"].(string)
		if !ok {
			log.Println("Invalid token data: username not found")
			http.Error(w, "Invalid token data", http.StatusUnauthorized)
			return
		}
	} else {
		// Если токена нет, установим имя как 'Гость'
		username = "Гость"
	}

	// Создаем WebSocket соединение
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading to WebSocket:", err)
		return
	}
	defer conn.Close()

	// Логируем подключение
	log.Println("WebSocket connection established")

	// Добавляем клиента в WebSocket-менеджер
	h.WebSocketManager.AddClient(conn)

	// Слушаем и обрабатываем сообщения
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading message: %v\n", err)
			break // Закрываем соединение при ошибке
		}

		// Только если пользователь авторизован, сохраняем сообщение в базе данных
		if token != "" {
			message := &domain.Message{
				UserID:    userIDInt,
				Username:  username,
				Content:   string(msg),
				Timestamp: time.Now(),
			}

			// Сохраняем сообщение в базе данных
			err = h.UseCase.SaveMessage(message)
			if err != nil {
				log.Println("Error saving message:", err)
				conn.WriteMessage(websocket.TextMessage, []byte("Error saving message"))
				continue
			}
		}

		// Отправляем сообщение всем подключенным клиентам
		h.WebSocketManager.BroadcastMessage([]byte(fmt.Sprintf(`{"username": "%s", "content": "%s"}`, username, string(msg))))
		log.Println("Message sent: ", string(msg))
	}

	// Убираем клиента после отключения
	h.WebSocketManager.RemoveClient(conn)
	log.Println("WebSocket connection closed")
}
