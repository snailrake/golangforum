package impl

import (
	"database/sql"
	"golangforum/internal/model"
)

type CommentRepository struct {
	db *sql.DB
}

func NewCommentRepository(db *sql.DB) *CommentRepository {
	return &CommentRepository{db: db}
}

func (r *CommentRepository) Create(c *model.Comment) error {
	_, err := r.db.Exec(
		"INSERT INTO comments (post_id, user_id, username, content, timestamp) VALUES ($1, $2, $3, $4, $5)",
		c.PostID, c.UserID, c.Username, c.Content, c.Timestamp,
	)
	return err
}

func (r *CommentRepository) GetByPost(postID int) ([]model.Comment, error) {
	rows, err := r.db.Query(
		"SELECT id, post_id, user_id, username, content, timestamp FROM comments WHERE post_id = $1",
		postID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []model.Comment
	for rows.Next() {
		var c model.Comment
		if err := rows.Scan(&c.ID, &c.PostID, &c.UserID, &c.Username, &c.Content, &c.Timestamp); err != nil {
			return nil, err
		}
		res = append(res, c)
	}
	return res, nil
}

func (r *CommentRepository) Delete(id int) error {
	_, err := r.db.Exec("DELETE FROM comments WHERE id = $1", id)
	return err
}
