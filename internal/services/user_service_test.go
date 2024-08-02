package services

import (
	"context"
	"testing"

	"gitlab.com/jkozhemiaka/web-layout/internal/constants"
	"gitlab.com/jkozhemiaka/web-layout/internal/models"
	mocks "gitlab.com/jkozhemiaka/web-layout/internal/repositories/mocks"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest"
)

func TestUserService_CreateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepoInterface(ctrl)
	mockLogger := zaptest.NewLogger(t).Sugar()
	userService := NewUserService(mockRepo, mockLogger)

	testUser := &models.User{Email: "test@example.com"}
	mockRepo.EXPECT().CreateUser(gomock.Any(), testUser).Return(testUser, nil)

	userId, err := userService.CreateUser(context.Background(), testUser)
	assert.NoError(t, err)
	assert.NotEqual(t, constants.EmptyString, userId)
}

func TestUserService_GetUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepoInterface(ctrl)
	mockLogger := zaptest.NewLogger(t).Sugar()
	userService := NewUserService(mockRepo, mockLogger)

	testUserID := "1"
	testUser := &models.User{ID: 1, Email: "test@example.com"}
	mockRepo.EXPECT().GetUser(gomock.Any(), testUserID).Return(testUser, nil)

	user, err := userService.GetUser(context.Background(), testUserID)
	assert.NoError(t, err)
	assert.Equal(t, testUser, user)
}

func TestUserService_DeleteUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepoInterface(ctrl)
	mockLogger := zaptest.NewLogger(t).Sugar()
	userService := NewUserService(mockRepo, mockLogger)

	testUserID := "1"
	testUser := &models.User{ID: 1, Email: "test@example.com"}
	mockRepo.EXPECT().DeleteUser(gomock.Any(), testUserID).Return(testUser, nil)

	user, err := userService.DeleteUser(context.Background(), testUserID)
	assert.NoError(t, err)
	assert.Equal(t, testUser, user)
}

func TestUserService_UpdateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepoInterface(ctrl)
	mockLogger := zaptest.NewLogger(t).Sugar()
	userService := NewUserService(mockRepo, mockLogger)

	testUserID := "1"
	testUser := &models.User{ID: 1, Email: "updated@example.com"}
	mockRepo.EXPECT().UpdateUser(gomock.Any(), testUserID, testUser).Return(testUser, nil)

	user, err := userService.UpdateUser(context.Background(), testUserID, testUser)
	assert.NoError(t, err)
	assert.Equal(t, testUser, user)
}

func TestUserService_ListUsers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepoInterface(ctrl)
	mockLogger := zaptest.NewLogger(t).Sugar()
	userService := NewUserService(mockRepo, mockLogger)

	testUsers := []models.User{
		{ID: 1, Email: "user1@example.com"},
		{ID: 2, Email: "user2@example.com"},
	}
	mockRepo.EXPECT().ListUsers(gomock.Any(), 1, 10).Return(testUsers, nil)

	users, err := userService.ListUsers(context.Background(), 1, 10)
	assert.NoError(t, err)
	assert.Equal(t, testUsers, users)
}

func TestUserService_CountUsers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepoInterface(ctrl)
	mockLogger := zaptest.NewLogger(t).Sugar()
	userService := NewUserService(mockRepo, mockLogger)

	mockRepo.EXPECT().CountUsers(gomock.Any()).Return(2, nil)

	count, err := userService.CountUsers(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, 2, count)
}

func TestUserService_GetUserByEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepoInterface(ctrl)
	mockLogger := zaptest.NewLogger(t).Sugar()
	userService := NewUserService(mockRepo, mockLogger)

	testEmail := "test@example.com"
	testUser := &models.User{ID: 1, Email: testEmail}
	mockRepo.EXPECT().GetUserByEmail(gomock.Any(), testEmail).Return(testUser, nil)

	user, err := userService.GetUserByEmail(context.Background(), testEmail)
	assert.NoError(t, err)
	assert.Equal(t, testUser, user)
}
