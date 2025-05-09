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

func TestPostUseCase_CRUD(t *testing.T) {
	ctx := context.Background()
	db, cleanup, err := utils.SetupPostgres(ctx, "db")
	if err != nil {
		t.Fatal(err)
	}
	defer cleanup()
	db.Exec(`INSERT INTO topics (id,title,description,created_at) VALUES (1,'t','',now()) ON CONFLICT DO NOTHING`)
	db.Exec(`TRUNCATE posts RESTART IDENTITY CASCADE`)

	r := impl.NewPostRepository(db)
	uc := usecase.NewPostUseCase(r)

	now := time.Now().Truncate(time.Second)
	p := &model.Post{TopicID: 1, Title: "Hello", Content: "World", UserID: 2}
	assert.NoError(t, uc.Create("alice", p))
	assert.Equal(t, "alice", p.Username)
	assert.WithinDuration(t, now, p.Timestamp, time.Minute)

	byTopic, err := uc.GetByTopic(1)
	assert.NoError(t, err)
	assert.Len(t, byTopic, 1)
	assert.Equal(t, "World", byTopic[0].Content)

	all, err := uc.GetAll()
	assert.NoError(t, err)
	assert.Len(t, all, 1)

	assert.NoError(t, uc.Delete(byTopic[0].ID))

	byTopic, err = uc.GetByTopic(1)
	assert.NoError(t, err)
	assert.Empty(t, byTopic)

	all, err = uc.GetAll()
	assert.NoError(t, err)
	assert.Empty(t, all)
}
