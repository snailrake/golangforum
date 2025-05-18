package usecase

import (
	"golangforum/internal/repository"
	"golangforum/internal/usecase"
	"time"

	"github.com/rs/zerolog/log"
	"golangforum/internal/model"
)

type CommentUseCase struct {
	repo repository.CommentRepository
}

func NewCommentUseCase(repo repository.CommentRepository) *CommentUseCase {
	log.Info().Msg("CommentUseCase initialized")
	return &CommentUseCase{repo: repo}
}

func (uc *CommentUseCase) Create(username string, c *model.Comment) error {
	log.Debug().Str("username", username).Msg("Creating comment")
	c.Username = username
	c.Timestamp = time.Now()
	if err := c.Validate(); err != nil {
		log.Warn().Err(err).Msg("Comment validation failed")
		return usecase.ErrInvalidCommentData
	}
	if err := uc.repo.Create(c); err != nil {
		log.Error().Err(err).Msg("Failed to save comment")
		return err
	}
	log.Info().
		Str("username", username).
		Time("timestamp", c.Timestamp).
		Msg("Comment created")
	return nil
}

func (uc *CommentUseCase) GetByPost(postID int) ([]model.Comment, error) {
	log.Debug().Int("postID", postID).Msg("Fetching comments for post")
	comments, err := uc.repo.GetByPost(postID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to fetch comments")
		return nil, err
	}
	log.Info().
		Int("postID", postID).
		Int("count", len(comments)).
		Msg("Comments fetched")
	return comments, nil
}

func (uc *CommentUseCase) Delete(id int) error {
	log.Debug().Int("id", id).Msg("Deleting comment")
	if err := uc.repo.Delete(id); err != nil {
		log.Error().Err(err).Msg("Failed to delete comment")
		if err.Error() == "not found" {
			return usecase.ErrCommentNotFound
		}
		return err
	}
	log.Info().Int("id", id).Msg("Comment deleted")
	return nil
}
