package usecase

import (
	"golangforum/internal/repository"
	"time"

	"github.com/rs/zerolog/log"
	"golangforum/internal/model"
)

type TopicUseCase struct {
	Repo repository.TopicRepository
}

func NewTopicUseCase(repo repository.TopicRepository) *TopicUseCase {
	log.Info().Msg("TopicUseCase initialized")
	return &TopicUseCase{Repo: repo}
}

func (uc *TopicUseCase) Create(title, description string) error {
	log.Debug().
		Str("title", title).
		Str("description", description).
		Msg("Creating topic")
	t := &model.Topic{
		Title:       title,
		Description: description,
		CreatedAt:   time.Now(),
	}
	if err := t.Validate(); err != nil {
		log.Warn().Err(err).Msg("Topic validation failed")
		return err
	}
	if err := uc.Repo.Create(t); err != nil {
		log.Error().Err(err).Msg("Failed to save topic")
		return err
	}
	log.Info().
		Str("title", title).
		Time("createdAt", t.CreatedAt).
		Msg("Topic created")
	return nil
}

func (uc *TopicUseCase) GetAll() ([]model.Topic, error) {
	log.Debug().Msg("Fetching all topics")
	topics, err := uc.Repo.GetAll()
	if err != nil {
		log.Error().Err(err).Msg("Failed to fetch topics")
		return nil, err
	}
	log.Info().
		Int("count", len(topics)).
		Msg("Topics fetched")
	return topics, nil
}

func (uc *TopicUseCase) Delete(id int) error {
	log.Debug().
		Int("id", id).
		Msg("Deleting topic")
	if err := uc.Repo.Delete(id); err != nil {
		log.Error().Err(err).Msg("Failed to delete topic")
		return err
	}
	log.Info().
		Int("id", id).
		Msg("Topic deleted")
	return nil
}
