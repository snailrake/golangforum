package usecase

import "golangforum/internal/model"

type TopicUseCase interface {
	Create(title, description string) error
	GetAll() ([]model.Topic, error)
	Delete(id int) error
}
