// repository/postgres/chat_repository.go
package postgres

import (
	"database/sql"
	"golangforum/internal/domain"
	"time"
)

type ChatRepository struct {
	DB *sql.DB
}

func NewChatRepository(db *sql.DB) *ChatRepository {
	return &ChatRepository{DB: db}
}

// SaveMessage сохраняет сообщение в базе данных
func (r *ChatRepository) SaveMessage(message *domain.Message) error {
	_, err := r.DB.Exec(
		"INSERT INTO messages (user_id, username, content, timestamp) VALUES ($1, $2, $3, $4)",
		message.UserID, message.Username, message.Content, message.Timestamp,
	)
	return err
}

// DeleteMessagesOlderThan удаляет сообщения старше заданного времени
func (r *ChatRepository) DeleteMessagesOlderThan(cutoffTime time.Time) error {
	_, err := r.DB.Exec("DELETE FROM messages WHERE timestamp < $1", cutoffTime)
	return err
}

// GetAllMessages возвращает все сообщения из базы данных
func (r *ChatRepository) GetAllMessages() ([]domain.Message, error) {
	rows, err := r.DB.Query("SELECT id, user_id, username, content, timestamp FROM messages")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []domain.Message
	for rows.Next() {
		var message domain.Message
		if err := rows.Scan(&message.ID, &message.UserID, &message.Username, &message.Content, &message.Timestamp); err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}

	return messages, nil
}
