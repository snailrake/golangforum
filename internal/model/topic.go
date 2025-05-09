package model

import (
	"errors"
	"strings"
	"time"
)

// Topic представляет собой тему форума
// @Description Структура темы с необходимыми полями для хранения данных о теме
type Topic struct {
	ID          int       `json:"id"`          // ID темы
	Title       string    `json:"title"`       // Заголовок темы
	Description string    `json:"description"` // Описание темы
	CreatedAt   time.Time `json:"created_at"`  // Дата и время создания темы
}

func (t *Topic) Validate() error {
	if strings.TrimSpace(t.Title) == "" {
		return errors.New("title cannot be empty")
	}
	if strings.TrimSpace(t.Description) == "" {
		return errors.New("description cannot be empty")
	}
	return nil
}
