package usecase

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"golangforum/internal/model"
	"golangforum/internal/repository/mocks"
	"golangforum/internal/usecase"
)

func TestPostUseCase_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockPostRepository(ctrl)

	post := &model.Post{
		ID:      1,
		TopicID: 1,
		Title:   "Valid Title",
		Content: "This is a post",
	}

	mockRepo.EXPECT().Create(gomock.Eq(post)).Return(nil).Times(1)

	uc := usecase.NewPostUseCase(mockRepo)
	err := uc.Create("testUser", post)

	assert.NoError(t, err)
	assert.Equal(t, "testUser", post.Username)
	assert.NotZero(t, post.Timestamp)
}

func TestPostUseCase_Create_ValidationError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockPostRepository(ctrl)

	post := &model.Post{
		ID:      1,
		TopicID: 1,
		Content: "",
	}

	uc := usecase.NewPostUseCase(mockRepo)
	err := uc.Create("testUser", post)

	assert.Error(t, err)
}

func TestPostUseCase_GetByTopic(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockPostRepository(ctrl)

	posts := []model.Post{
		{ID: 1, TopicID: 1, Content: "First post", Timestamp: time.Now()},
		{ID: 2, TopicID: 1, Content: "Second post", Timestamp: time.Now()},
	}

	mockRepo.EXPECT().GetByTopic(1).Return(posts, nil).Times(1)

	uc := usecase.NewPostUseCase(mockRepo)
	result, err := uc.GetByTopic(1)

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, posts, result)
}

func TestPostUseCase_GetByTopic_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockPostRepository(ctrl)

	mockRepo.EXPECT().GetByTopic(1).Return(nil, errors.New("database error")).Times(1)

	uc := usecase.NewPostUseCase(mockRepo)
	result, err := uc.GetByTopic(1)

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestPostUseCase_GetAll(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockPostRepository(ctrl)

	posts := []model.Post{
		{ID: 1, TopicID: 1, Content: "First post", Timestamp: time.Now()},
		{ID: 2, TopicID: 1, Content: "Second post", Timestamp: time.Now()},
	}

	mockRepo.EXPECT().GetAll().Return(posts, nil).Times(1)

	uc := usecase.NewPostUseCase(mockRepo)
	result, err := uc.GetAll()

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, posts, result)
}

func TestPostUseCase_GetAll_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockPostRepository(ctrl)

	mockRepo.EXPECT().GetAll().Return(nil, errors.New("database error")).Times(1)

	uc := usecase.NewPostUseCase(mockRepo)
	result, err := uc.GetAll()

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestPostUseCase_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockPostRepository(ctrl)

	mockRepo.EXPECT().Delete(1).Return(nil).Times(1)

	uc := usecase.NewPostUseCase(mockRepo)
	err := uc.Delete(1)

	assert.NoError(t, err)
}

func TestPostUseCase_Delete_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockPostRepository(ctrl)

	mockRepo.EXPECT().Delete(1).Return(errors.New("delete error")).Times(1)

	uc := usecase.NewPostUseCase(mockRepo)
	err := uc.Delete(1)

	assert.Error(t, err)
}
