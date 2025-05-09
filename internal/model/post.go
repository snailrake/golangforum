package model

import (
	"errors"
	"strings"
	"time"
)

// Post представляет собой пост на форуме
// @Description Структура поста с необходимыми полями для хранения данных о посте
type Post struct {
	ID        int       `json:"id"`        // ID поста
	TopicID   int       `json:"topic_id"`  // ID темы, к которой относится пост
	Title     string    `json:"title"`     // Заголовок поста
	Content   string    `json:"content"`   // Текст поста
	UserID    int       `json:"user_id"`   // ID пользователя, создавшего пост
	Username  string    `json:"username"`  // Имя пользователя
	Timestamp time.Time `json:"timestamp"` // Время создания поста
}

func (p *Post) Validate() error {
	if p.TopicID <= 0 {
		return errors.New("topic_id must be positive")
	}
	if strings.TrimSpace(p.Title) == "" {
		return errors.New("title cannot be empty")
	}
	if strings.TrimSpace(p.Content) == "" {
		return errors.New("content cannot be empty")
	}
	return nil
}
