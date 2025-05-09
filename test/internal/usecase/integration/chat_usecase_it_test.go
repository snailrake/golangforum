package usecase

import (
	"context"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"golangforum/internal/model"
	"golangforum/internal/repository/impl"
	"golangforum/internal/usecase"
	"golangforum/test/utils"
)

func TestChatRepository_CRUD(t *testing.T) {
	ctx := context.Background()
	db, cleanup, err := utils.SetupPostgres(ctx, "db")
	if err != nil {
		t.Fatal(err)
	}
	defer cleanup()
	_, _ = db.Exec("TRUNCATE messages RESTART IDENTITY CASCADE")

	repo := impl.NewChatRepository(db)
	now := time.Now()
	assert.NoError(t, repo.SaveMessage(&model.Message{UserID: 1, Username: "u1", Content: "hello", Timestamp: now}))
	msgs, err := repo.GetAllMessages()
	assert.NoError(t, err)
	assert.Len(t, msgs, 1)
	assert.Equal(t, "hello", msgs[0].Content)

	assert.NoError(t, repo.SaveMessage(&model.Message{UserID: 2, Username: "u2", Content: "old", Timestamp: now.Add(-48 * time.Hour)}))
	assert.NoError(t, repo.DeleteMessagesOlderThan(now.Add(-24*time.Hour)))
	msgs, err = repo.GetAllMessages()
	assert.NoError(t, err)
	assert.Len(t, msgs, 1)
	assert.Equal(t, "hello", msgs[0].Content)
}

func TestChatUseCase_WithPostgres(t *testing.T) {
	ctx := context.Background()
	db, cleanup, err := utils.SetupPostgres(ctx, "db")
	if err != nil {
		t.Fatal(err)
	}
	defer cleanup()
	_, _ = db.Exec("TRUNCATE messages RESTART IDENTITY CASCADE")

	repo := impl.NewChatRepository(db)
	uc := usecase.NewChatUseCase(repo)
	msgs, err := uc.GetAllMessages()
	assert.NoError(t, err)
	assert.Empty(t, msgs)

	now := time.Now()
	repo.SaveMessage(&model.Message{UserID: 1, Username: "u1", Content: "first", Timestamp: now.Add(-time.Minute)})
	repo.SaveMessage(&model.Message{UserID: 2, Username: "u2", Content: "second", Timestamp: now})

	msgs, err = uc.GetAllMessages()
	assert.NoError(t, err)
	assert.Len(t, msgs, 2)
	assert.True(t, msgs[0].Timestamp.Before(msgs[1].Timestamp))
}
