package postgres

import (
	"database/sql"
	"golangforum/internal/domain"
)

type PostRepository struct {
	DB *sql.DB
}

func NewPostRepository(db *sql.DB) *PostRepository {
	return &PostRepository{DB: db}
}

func (r *PostRepository) Create(post *domain.Post) error {
	_, err := r.DB.Exec(
		"INSERT INTO posts (topic_id, title, content, user_id, username, timestamp) VALUES ($1, $2, $3, $4, $5, $6)",
		post.TopicID, post.Title, post.Content, post.UserID, post.Username, post.Timestamp,
	)
	return err
}

func (r *PostRepository) GetByTopic(topicID int) ([]domain.Post, error) {
	rows, err := r.DB.Query(
		"SELECT id, topic_id, title, content, user_id, username, timestamp FROM posts WHERE topic_id = $1",
		topicID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []domain.Post
	for rows.Next() {
		var p domain.Post
		if err := rows.Scan(&p.ID, &p.TopicID, &p.Title, &p.Content, &p.UserID, &p.Username, &p.Timestamp); err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}
	return posts, nil
}
