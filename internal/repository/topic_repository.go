package repository

import (
	"golangforum/internal/model"
)

type TopicRepository interface {
	Create(topic *model.Topic) error
	GetAll() ([]model.Topic, error)
	Delete(id int) error
}
