package usecase

import (
	"golangforum/internal/repository"
	"golangforum/internal/usecase"
	"time"

	"github.com/rs/zerolog/log"
	"golangforum/internal/model"
)

type PostUseCase struct {
	Repo repository.PostRepository
}

func NewPostUseCase(repo repository.PostRepository) *PostUseCase {
	log.Info().Msg("PostUseCase initialized")
	return &PostUseCase{Repo: repo}
}

func (uc *PostUseCase) Create(username string, post *model.Post) error {
	log.Debug().
		Str("username", username).
		Msg("Creating post")
	post.Username = username
	post.Timestamp = time.Now()
	if err := post.Validate(); err != nil {
		log.Warn().Err(err).Msg("Post validation failed")
		return usecase.ErrInvalidPostData
	}
	if err := uc.Repo.Create(post); err != nil {
		log.Error().Err(err).Msg("Failed to save post")
		return err
	}
	log.Info().
		Str("username", username).
		Time("timestamp", post.Timestamp).
		Msg("Post created")
	return nil
}

func (uc *PostUseCase) GetByTopic(topicID int) ([]model.Post, error) {
	log.Debug().
		Int("topicID", topicID).
		Msg("Fetching posts by topic")
	posts, err := uc.Repo.GetByTopic(topicID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to fetch posts by topic")
		return nil, err
	}
	log.Info().
		Int("topicID", topicID).
		Int("count", len(posts)).
		Msg("Posts fetched by topic")
	return posts, nil
}

func (uc *PostUseCase) GetAll() ([]model.Post, error) {
	log.Debug().Msg("Fetching all posts")
	posts, err := uc.Repo.GetAll()
	if err != nil {
		log.Error().Err(err).Msg("Failed to fetch all posts")
		return nil, err
	}
	log.Info().
		Int("count", len(posts)).
		Msg("All posts fetched")
	return posts, nil
}

func (uc *PostUseCase) Delete(id int) error {
	log.Debug().
		Int("id", id).
		Msg("Deleting post")
	if err := uc.Repo.Delete(id); err != nil {
		log.Error().Err(err).Msg("Failed to delete post")
		if err.Error() == "not found" {
			return usecase.ErrPostNotFound
		}
		return err
	}
	log.Info().
		Int("id", id).
		Msg("Post deleted")
	return nil
}
