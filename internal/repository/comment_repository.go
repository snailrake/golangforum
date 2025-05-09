package repository

import (
	"golangforum/internal/model"
)

type CommentRepository interface {
	Create(c *model.Comment) error
	GetByPost(postID int) ([]model.Comment, error)
	Delete(id int) error
}
