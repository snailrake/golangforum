package postgres

import (
	"database/sql"
	"golangforum/internal/domain"
)

type CommentRepository struct {
	DB *sql.DB
}

func NewCommentRepository(db *sql.DB) *CommentRepository {
	return &CommentRepository{DB: db}
}

func (r *CommentRepository) Create(comment *domain.Comment) error {
	_, err := r.DB.Exec(
		"INSERT INTO comments (post_id, user_id, username, content, timestamp) VALUES ($1, $2, $3, $4, $5)",
		comment.PostID, comment.UserID, comment.Username, comment.Content, comment.Timestamp,
	)
	return err
}

func (r *CommentRepository) GetByPost(postID int) ([]domain.Comment, error) {
	rows, err := r.DB.Query(
		"SELECT id, post_id, user_id, username, content, timestamp FROM comments WHERE post_id = $1",
		postID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []domain.Comment
	for rows.Next() {
		var c domain.Comment
		if err := rows.Scan(&c.ID, &c.PostID, &c.UserID, &c.Username, &c.Content, &c.Timestamp); err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}
	return comments, nil
}
