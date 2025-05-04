package usecase

import (
	"golangforum/internal/domain"
	"golangforum/internal/repository/postgres"
	"time"
)

type CommentUseCase struct {
	Repo *postgres.CommentRepository
}

func NewCommentUseCase(repo *postgres.CommentRepository) *CommentUseCase {
	return &CommentUseCase{Repo: repo}
}

func (uc *CommentUseCase) Create(comment *domain.Comment) error {
	comment.Timestamp = time.Now()
	return uc.Repo.Create(comment)
}

func (uc *CommentUseCase) GetByPost(postID int) ([]domain.Comment, error) {
	return uc.Repo.GetByPost(postID)
}
