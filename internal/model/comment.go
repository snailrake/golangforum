package model

import (
	"errors"
	"strings"
	"time"
)

// Comment представляет собой комментарий на посте
// @Description Структура комментария с необходимыми полями для хранения данных о комментарии
type Comment struct {
	ID        int       `json:"id"`        // ID комментария
	PostID    int       `json:"post_id"`   // ID поста, к которому относится комментарий
	UserID    int       `json:"user_id"`   // ID пользователя, оставившего комментарий
	Username  string    `json:"username"`  // Имя пользователя
	Content   string    `json:"content"`   // Текст комментария
	Timestamp time.Time `json:"timestamp"` // Время создания комментария
}

func (c *Comment) Validate() error {
	if c.PostID <= 0 {
		return errors.New("post_id must be positive")
	}
	if strings.TrimSpace(c.Content) == "" {
		return errors.New("content cannot be empty")
	}
	return nil
}
