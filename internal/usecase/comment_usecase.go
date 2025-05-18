package usecase

import "golangforum/internal/model"

type CommentUseCase interface {
	Create(username string, c *model.Comment) error
	GetByPost(postID int) ([]model.Comment, error)
	Delete(id int) error
}
