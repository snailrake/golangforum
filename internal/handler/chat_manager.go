package handler

import (
	"github.com/gorilla/websocket"
	"log"
)

type WebSocketManager struct {
	clients map[*websocket.Conn]bool
}

func NewWebSocketManager() *WebSocketManager {
	return &WebSocketManager{
		clients: make(map[*websocket.Conn]bool),
	}
}

// AddClient добавляет нового клиента в список
func (wm *WebSocketManager) AddClient(conn *websocket.Conn) {
	wm.clients[conn] = true
}

// RemoveClient удаляет клиента из списка
func (wm *WebSocketManager) RemoveClient(conn *websocket.Conn) {
	delete(wm.clients, conn)
}

// BroadcastMessage рассылает сообщение всем клиентам
func (wm *WebSocketManager) BroadcastMessage(message []byte) {
	for client := range wm.clients {
		if err := client.WriteMessage(websocket.TextMessage, message); err != nil {
			log.Println("Error broadcasting message:", err)
			client.Close()
			wm.RemoveClient(client)
		}
	}
}
