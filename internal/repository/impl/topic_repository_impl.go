package impl

import (
	"database/sql"

	"golangforum/internal/model"
)

type TopicRepository struct {
	DB *sql.DB
}

func NewTopicRepository(db *sql.DB) *TopicRepository {
	return &TopicRepository{DB: db}
}

func (r *TopicRepository) Create(topic *model.Topic) error {
	_, err := r.DB.Exec(
		"INSERT INTO topics (title, description, created_at) VALUES ($1, $2, $3)",
		topic.Title, topic.Description, topic.CreatedAt,
	)
	return err
}

func (r *TopicRepository) GetAll() ([]model.Topic, error) {
	rows, err := r.DB.Query(
		"SELECT id, title, description, created_at FROM topics",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var topics []model.Topic
	for rows.Next() {
		var t model.Topic
		if err := rows.Scan(&t.ID, &t.Title, &t.Description, &t.CreatedAt); err != nil {
			return nil, err
		}
		topics = append(topics, t)
	}
	return topics, rows.Err()
}

func (r *TopicRepository) Delete(id int) error {
	_, err := r.DB.Exec("DELETE FROM topics WHERE id = $1", id)
	return err
}
