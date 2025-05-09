package usecase

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"golangforum/internal/repository/impl"
	"golangforum/internal/usecase"
	"golangforum/test/utils"
)

func TestTopicUseCase_CRUD(t *testing.T) {
	ctx := context.Background()
	db, terminate, err := utils.SetupPostgres(ctx, "db")
	if err != nil {
		t.Fatalf("postgres setup: %v", err)
	}
	defer terminate()

	_, _ = db.Exec("TRUNCATE topics RESTART IDENTITY CASCADE")

	repo := impl.NewTopicRepository(db)
	uc := usecase.NewTopicUseCase(repo)

	assert.NoError(t, uc.Create("Title1", "Desc1"))

	topics, err := uc.GetAll()
	assert.NoError(t, err)
	assert.Len(t, topics, 1)
	topic := topics[0]
	assert.Equal(t, "Title1", topic.Title)
	assert.Equal(t, "Desc1", topic.Description)
	assert.False(t, topic.CreatedAt.IsZero())

	assert.NoError(t, uc.Delete(topic.ID))

	topics, err = uc.GetAll()
	assert.NoError(t, err)
	assert.Empty(t, topics)
}
