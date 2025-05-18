package usecase

import "golangforum/internal/model"

type PostUseCase interface {
	Create(username string, post *model.Post) error
	GetByTopic(topicID int) ([]model.Post, error)
	GetAll() ([]model.Post, error)
	Delete(id int) error
}
