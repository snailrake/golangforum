package model

import (
	"errors"
	"strings"
	"time"
)

// Message представляет собой сообщение от пользователя
// @Description Структура сообщения с необходимыми полями для хранения данных о сообщении
type Message struct {
	ID        int       `json:"id"`        // ID сообщения
	UserID    int       `json:"user_id"`   // ID пользователя, отправившего сообщение
	Username  string    `json:"username"`  // Имя пользователя
	Content   string    `json:"content"`   // Текст сообщения
	Timestamp time.Time `json:"timestamp"` // Время отправки сообщения
}

func (m *Message) Validate() error {
	if strings.TrimSpace(m.Username) == "" {
		return errors.New("username cannot be empty")
	}
	if strings.TrimSpace(m.Content) == "" {
		return errors.New("content cannot be empty")
	}
	return nil
}
