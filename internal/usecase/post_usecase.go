package usecase

import (
	"golangforum/internal/domain"
	"golangforum/internal/repository/postgres"
	"time"
)

type PostUseCase struct {
	Repo *postgres.PostRepository
}

func NewPostUseCase(repo *postgres.PostRepository) *PostUseCase {
	return &PostUseCase{Repo: repo}
}

func (uc *PostUseCase) Create(post *domain.Post) error {
	post.Timestamp = time.Now()
	return uc.Repo.Create(post)
}

func (uc *PostUseCase) GetByTopic(topicID int) ([]domain.Post, error) {
	return uc.Repo.GetByTopic(topicID)
}

func (uc *PostUseCase) GetAll() ([]domain.Post, error) {
	return uc.Repo.GetAll()
}
