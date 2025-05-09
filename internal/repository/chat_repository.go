package repository

import (
	"golangforum/internal/model"
	"time"
)

type ChatRepository interface {
	GetAllMessages() ([]model.Message, error)
	SaveMessage(*model.Message) error
	DeleteMessagesOlderThan(time.Time) error
}
