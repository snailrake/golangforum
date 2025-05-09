package impl

import (
	"database/sql"
	"golangforum/internal/model"
	"time"
)

type ChatRepositoryImpl struct {
	DB *sql.DB
}

func NewChatRepository(db *sql.DB) *ChatRepositoryImpl {
	return &ChatRepositoryImpl{DB: db}
}

func (r *ChatRepositoryImpl) SaveMessage(m *model.Message) error {
	_, err := r.DB.Exec(
		"INSERT INTO messages (user_id, username, content, timestamp) VALUES ($1, $2, $3, $4)",
		m.UserID, m.Username, m.Content, m.Timestamp,
	)
	return err
}

func (r *ChatRepositoryImpl) DeleteMessagesOlderThan(t time.Time) error {
	_, err := r.DB.Exec("DELETE FROM messages WHERE timestamp < $1", t)
	return err
}

func (r *ChatRepositoryImpl) GetAllMessages() ([]model.Message, error) {
	rows, err := r.DB.Query("SELECT id, user_id, username, content, timestamp FROM messages")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var msgs []model.Message
	for rows.Next() {
		var m model.Message
		if err := rows.Scan(&m.ID, &m.UserID, &m.Username, &m.Content, &m.Timestamp); err != nil {
			return nil, err
		}
		msgs = append(msgs, m)
	}

	return msgs, nil
}
