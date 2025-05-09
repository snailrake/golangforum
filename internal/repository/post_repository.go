package repository

import (
	"golangforum/internal/model"
)

type PostRepository interface {
	Create(post *model.Post) error
	GetByTopic(topicID int) ([]model.Post, error)
	GetAll() ([]model.Post, error)
	Delete(id int) error
}
