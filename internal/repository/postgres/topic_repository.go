package postgres

import (
	"database/sql"
	"golangforum/internal/domain"
)

type TopicRepository struct {
	DB *sql.DB
}

func NewTopicRepository(db *sql.DB) *TopicRepository {
	return &TopicRepository{DB: db}
}

func (r *TopicRepository) Create(topic *domain.Topic) error {
	_, err := r.DB.Exec(
		"INSERT INTO topics (title, description, created_at) VALUES ($1, $2, $3)",
		topic.Title, topic.Description, topic.CreatedAt,
	)
	return err
}

func (r *TopicRepository) GetAll() ([]domain.Topic, error) {
	rows, err := r.DB.Query(
		"SELECT id, title, description, created_at FROM topics",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var topics []domain.Topic
	for rows.Next() {
		var t domain.Topic
		if err := rows.Scan(&t.ID, &t.Title, &t.Description, &t.CreatedAt); err != nil {
			return nil, err
		}
		topics = append(topics, t)
	}
	return topics, nil
}
