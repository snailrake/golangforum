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

func TestCommentUseCase_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockCommentRepository(ctrl)
	comment := &model.Comment{ID: 1, PostID: 1, Content: "This is a comment"}

	mockRepo.EXPECT().Create(comment).Return(nil).Times(1)

	uc := usecase.NewCommentUseCase(mockRepo)
	err := uc.Create("testUser", comment)

	assert.NoError(t, err)
	assert.Equal(t, "testUser", comment.Username)
	assert.NotZero(t, comment.Timestamp)
}

func TestCommentUseCase_Create_ValidationError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockCommentRepository(ctrl)
	comment := &model.Comment{ID: 1, PostID: 1, Content: ""}

	uc := usecase.NewCommentUseCase(mockRepo)
	err := uc.Create("testUser", comment)

	assert.Error(t, err)
}

func TestCommentUseCase_GetByPost(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockCommentRepository(ctrl)
	comments := []model.Comment{
		{ID: 1, PostID: 1, Content: "First comment", Timestamp: time.Now()},
		{ID: 2, PostID: 1, Content: "Second comment", Timestamp: time.Now()},
	}

	mockRepo.EXPECT().GetByPost(1).Return(comments, nil).Times(1)

	uc := usecase.NewCommentUseCase(mockRepo)
	result, err := uc.GetByPost(1)

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, comments, result)
}

func TestCommentUseCase_GetByPost_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockCommentRepository(ctrl)
	mockRepo.EXPECT().GetByPost(1).Return(nil, errors.New("database error")).Times(1)

	uc := usecase.NewCommentUseCase(mockRepo)
	result, err := uc.GetByPost(1)

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestCommentUseCase_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockCommentRepository(ctrl)
	mockRepo.EXPECT().Delete(1).Return(nil).Times(1)

	uc := usecase.NewCommentUseCase(mockRepo)
	err := uc.Delete(1)

	assert.NoError(t, err)
}

func TestCommentUseCase_Delete_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockCommentRepository(ctrl)
	mockRepo.EXPECT().Delete(1).Return(errors.New("delete error")).Times(1)

	uc := usecase.NewCommentUseCase(mockRepo)
	err := uc.Delete(1)

	assert.Error(t, err)
}
