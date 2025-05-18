package usecase

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"golangforum/internal/model"
	"golangforum/internal/repository/mocks"
)

func TestChatUseCase_GetAllMessages(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockChatRepository(ctrl)
	msgs := []model.Message{
		{Timestamp: time.Now().Add(-2 * time.Hour)},
		{Timestamp: time.Now().Add(-1 * time.Hour)},
	}
	mockRepo.EXPECT().GetAllMessages().Return(msgs, nil)

	uc := NewChatUseCase(mockRepo)
	result, err := uc.GetAllMessages()
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.True(t, result[0].Timestamp.Before(result[1].Timestamp))
}

func TestChatUseCase_GetAllMessages_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockChatRepository(ctrl)
	mockRepo.EXPECT().GetAllMessages().Return(nil, errors.New("fail"))

	uc := NewChatUseCase(mockRepo)
	result, err := uc.GetAllMessages()
	assert.Nil(t, result)
	assert.Error(t, err)
}

func TestChatUseCase_HandleConnection_ValidMessage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockChatRepository(ctrl)
	mockRepo.EXPECT().SaveMessage(gomock.Any()).Return(nil).Times(1)
	mockRepo.EXPECT().DeleteMessagesOlderThan(gomock.Any()).Return(nil).Times(1)

	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			t.Fatal(err)
		}
		uc := NewChatUseCase(mockRepo)
		uc.HandleConnection(conn, "testUser", 1)
	})
	server := httptest.NewServer(handler)
	defer server.Close()

	wsURL := "ws" + server.URL[4:]
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatal("Dial:", err)
	}
	defer conn.Close()

	err = conn.WriteMessage(websocket.TextMessage, []byte("Hello"))
	assert.NoError(t, err)

	_, msg, err := conn.ReadMessage()
	assert.NoError(t, err)
	assert.Contains(t, string(msg), "Hello")
}
