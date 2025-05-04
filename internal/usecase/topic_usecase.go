package usecase

import (
	"golangforum/internal/domain"
	"golangforum/internal/repository/postgres"
	"time"
)

type TopicUseCase struct {
	Repo *postgres.TopicRepository
}

func NewTopicUseCase(repo *postgres.TopicRepository) *TopicUseCase {
	return &TopicUseCase{Repo: repo}
}

func (uc *TopicUseCase) Create(topic *domain.Topic) error {
	topic.CreatedAt = time.Now()
	return uc.Repo.Create(topic)
}

func (uc *TopicUseCase) GetAll() ([]domain.Topic, error) {
	return uc.Repo.GetAll()
}
