package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"golangforum/internal/model"
	"golangforum/internal/repository/impl"
	"golangforum/internal/usecase"
	"golangforum/test/utils"
)

func TestCommentUseCase_CRUD(t *testing.T) {
	db, terminate, err := utils.SetupPostgres(context.Background(), "db")
	if err != nil {
		t.Fatalf("postgres setup: %v", err)
	}
	defer terminate()

	repo := impl.NewCommentRepository(db)
	uc := usecase.NewCommentUseCase(repo)

	now := time.Now().Truncate(time.Second)
	c := &model.Comment{PostID: 1, UserID: 2, Content: "nice"}
	err = uc.Create("bob", c)
	assert.NoError(t, err)
	assert.Equal(t, "bob", c.Username)
	assert.WithinDuration(t, now, c.Timestamp, time.Minute)

	comments, err := uc.GetByPost(1)
	assert.NoError(t, err)
	assert.Len(t, comments, 1)
	assert.Equal(t, c.Content, comments[0].Content)
	assert.Equal(t, c.UserID, comments[0].UserID)

	err = uc.Delete(comments[0].ID)
	assert.NoError(t, err)

	comments, err = uc.GetByPost(1)
	assert.NoError(t, err)
	assert.Empty(t, comments)
}
