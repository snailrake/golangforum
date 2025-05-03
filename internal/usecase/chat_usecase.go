// usecase/chat_usecase.go
package usecase

import (
	"golangforum/internal/domain"
	"golangforum/internal/repository/postgres"
	"time"
)

type ChatUseCase struct {
	Repo *postgres.ChatRepository
}

func NewChatUseCase(repo *postgres.ChatRepository) *ChatUseCase {
	return &ChatUseCase{Repo: repo}
}

// SaveMessage сохраняет сообщение и удаляет старые сообщения
func (uc *ChatUseCase) SaveMessage(message *domain.Message) error {
	// Сохраняем сообщение в базе данных
	if err := uc.Repo.SaveMessage(message); err != nil {
		return err
	}

	// Удаляем старые сообщения на основе времени хранения
	return uc.DeleteOldMessages()
}

// DeleteOldMessages удаляет сообщения старше заданного времени
func (uc *ChatUseCase) DeleteOldMessages() error {
	// Получаем период хранения сообщений из .env
	period, err := time.ParseDuration("24h")
	if err != nil {
		return err
	}

	// Удаляем сообщения старше периода хранения
	return uc.Repo.DeleteMessagesOlderThan(time.Now().Add(-period))
}

// GetAllMessages получает все сообщения из базы данных
func (uc *ChatUseCase) GetAllMessages() ([]domain.Message, error) {
	return uc.Repo.GetAllMessages()
}
