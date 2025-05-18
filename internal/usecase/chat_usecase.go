package usecase

import (
	"github.com/gorilla/websocket"
	"golangforum/internal/model"
)

type ChatUseCase interface {
	GetAllMessages() ([]model.Message, error)
	HandleConnection(conn *websocket.Conn, user string, id int, clients map[*websocket.Conn]struct{})
}
