package usecase

import (
	"encoding/json"
	"golangforum/internal/repository"
	"sort"
	"time"

	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
	"golangforum/internal/model"
)

type ChatUseCase struct {
	repo    repository.ChatRepository
	clients map[*websocket.Conn]struct{}
}

func NewChatUseCase(repo repository.ChatRepository) *ChatUseCase {
	log.Info().Msg("Initializing ChatUseCase")
	return &ChatUseCase{repo: repo, clients: make(map[*websocket.Conn]struct{})}
}

func (uc *ChatUseCase) GetAllMessages() ([]model.Message, error) {
	log.Info().Msg("GetAllMessages called")
	msgs, err := uc.repo.GetAllMessages()
	if err != nil {
		log.Error().Err(err).Msg("Error fetching messages from repository")
		return nil, err
	}
	sort.Slice(msgs, func(i, j int) bool {
		return msgs[i].Timestamp.Before(msgs[j].Timestamp)
	})
	log.Info().Int("count", len(msgs)).Msg("Messages sorted by timestamp")
	return msgs, nil
}

func (uc *ChatUseCase) HandleConnection(conn *websocket.Conn, user string, id int) {
	log.Info().
		Str("user", user).
		Int("userID", id).
		Msg("Handling new WebSocket connection")
	uc.clients[conn] = struct{}{}
	defer func() {
		log.Info().
			Str("user", user).
			Int("userID", id).
			Msg("Closing WebSocket connection")
		delete(uc.clients, conn)
		conn.Close()
	}()

	for {
		_, data, err := conn.ReadMessage()
		if err != nil {
			log.Warn().Err(err).Msg("ReadMessage error, terminating connection loop")
			return
		}
		if id == 0 {
			log.Warn().
				Str("user", user).
				Int("userID", id).
				Msg("Received message from invalid user ID, skipping")
			continue
		}
		text := string(data)
		log.Debug().
			Str("user", user).
			Int("userID", id).
			Str("message", text).
			Msg("Received message")

		m := &model.Message{UserID: id, Username: user, Content: text, Timestamp: time.Now()}
		if err := m.Validate(); err != nil {
			log.Warn().
				Str("user", user).
				Int("userID", id).
				Msg("Message validation failed")
		} else {
			if err := uc.repo.SaveMessage(m); err != nil {
				log.Error().Err(err).Msg("Failed to save message")
			} else {
				log.Info().
					Str("user", user).
					Int("userID", id).
					Msg("Message saved to repository")
			}
			if err := uc.repo.DeleteMessagesOlderThan(time.Now().Add(-24 * time.Hour)); err != nil {
				log.Error().Err(err).Msg("Failed to delete old messages")
			} else {
				log.Debug().Msg("Old messages cleanup completed")
			}
		}

		out, _ := json.Marshal(struct {
			Username string `json:"username"`
			Content  string `json:"content"`
		}{user, text})

		for c := range uc.clients {
			if err := c.WriteMessage(websocket.TextMessage, out); err != nil {
				log.Warn().Err(err).Msg("Error writing to client, removing connection")
				c.Close()
				delete(uc.clients, c)
			} else {
				log.Debug().Msg("Message broadcast to client")
			}
		}
	}
}
