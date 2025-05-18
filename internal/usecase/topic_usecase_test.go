package usecase

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"golangforum/internal/model"
	"golangforum/internal/repository/mocks"
)

func TestTopicUseCase_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockTopicRepository(ctrl)

	mockRepo.EXPECT().Create(gomock.Any()).Return(nil).Times(1)

	uc := NewTopicUseCase(mockRepo)
	err := uc.Create("Valid Topic", "This is a valid topic")

	assert.NoError(t, err)
}

func TestTopicUseCase_Create_ValidationError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockTopicRepository(ctrl)

	uc := NewTopicUseCase(mockRepo)
	err := uc.Create("", "This is a topic without a title")

	assert.Error(t, err)
}

func TestTopicUseCase_GetAll(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockTopicRepository(ctrl)
	topics := []model.Topic{
		{ID: 1, Title: "First Topic", Description: "This is the first topic", CreatedAt: time.Now()},
		{ID: 2, Title: "Second Topic", Description: "This is the second topic", CreatedAt: time.Now()},
	}

	mockRepo.EXPECT().GetAll().Return(topics, nil).Times(1)

	uc := NewTopicUseCase(mockRepo)
	result, err := uc.GetAll()

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, topics, result)
}

func TestTopicUseCase_GetAll_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockTopicRepository(ctrl)

	mockRepo.EXPECT().GetAll().Return(nil, errors.New("database error")).Times(1)

	uc := NewTopicUseCase(mockRepo)
	result, err := uc.GetAll()

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestTopicUseCase_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockTopicRepository(ctrl)

	mockRepo.EXPECT().Delete(1).Return(nil).Times(1)

	uc := NewTopicUseCase(mockRepo)
	err := uc.Delete(1)

	assert.NoError(t, err)
}

func TestTopicUseCase_Delete_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockTopicRepository(ctrl)

	mockRepo.EXPECT().Delete(1).Return(errors.New("delete error")).Times(1)

	uc := NewTopicUseCase(mockRepo)
	err := uc.Delete(1)

	assert.Error(t, err)
}
