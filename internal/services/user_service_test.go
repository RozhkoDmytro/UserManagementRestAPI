package services

import (
	"context"
	"errors"
	"strconv"
	"testing"
	"time"

	"gitlab.com/jkozhemiaka/web-layout/internal/apperrors"
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

func TestUserService_Vote(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepoInterface(ctrl)
	mockLogger := zaptest.NewLogger(t).Sugar()
	userService := NewUserService(mockRepo, mockLogger)

	testVote := &models.Vote{UserID: 1, ProfileID: 2, Value: 1}
	testUser := &models.User{ID: 1, UpdatedAt: time.Now().Add(-2 * time.Hour)}

	// Return the user and nil for error
	mockRepo.EXPECT().GetUserByID(gomock.Any(), testVote.UserID).Return(testUser, nil)
	mockRepo.EXPECT().GetVote(gomock.Any(), testVote.UserID, testVote.ProfileID).Return(nil, nil)
	mockRepo.EXPECT().CreateVote(gomock.Any(), testVote).Return(testVote, nil)

	voteID, err := userService.Vote(context.Background(), testVote)
	assert.NoError(t, err)
	assert.Equal(t, strconv.Itoa(int(testVote.ID)), voteID)
}

func TestUserService_Vote_CooldownError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepoInterface(ctrl)
	mockLogger := zaptest.NewLogger(t).Sugar()
	userService := NewUserService(mockRepo, mockLogger)

	testVote := &models.Vote{UserID: 1, ProfileID: 2, Value: 1}
	testUser := &models.User{ID: 1, UpdatedAt: time.Now().Add(-30 * time.Minute)} // Time within cooldown period

	// Set expectations
	mockRepo.EXPECT().GetUserByID(gomock.Any(), testVote.UserID).Return(testUser, nil)

	_, err := userService.Vote(context.Background(), testVote)
	assert.Error(t, err)
	// Type assertion to *apperrors.AppError
	appErr, ok := err.(*apperrors.AppError)
	assert.True(t, ok)
	assert.Equal(t, apperrors.VoteCooldownErr.Code, appErr.Code)
}

func TestUserService_Vote_UpdateSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepoInterface(ctrl)
	mockLogger := zaptest.NewLogger(t).Sugar()
	userService := NewUserService(mockRepo, mockLogger)

	testVote := &models.Vote{UserID: 1, ProfileID: 2, Value: 1}
	existingVote := &models.Vote{ID: 10, UserID: 1, ProfileID: 2, Value: 0}
	testUser := &models.User{ID: 1, UpdatedAt: time.Now().Add(-2 * time.Hour)} // Time outside cooldown period

	// Set expectations
	mockRepo.EXPECT().GetUserByID(gomock.Any(), testVote.UserID).Return(testUser, nil)
	mockRepo.EXPECT().GetVote(gomock.Any(), testVote.UserID, testVote.ProfileID).Return(existingVote, nil)
	mockRepo.EXPECT().UpdateVote(gomock.Any(), existingVote).Return(existingVote, nil)

	voteID, err := userService.Vote(context.Background(), testVote)
	assert.NoError(t, err)
	assert.Equal(t, strconv.Itoa(int(existingVote.ID)), voteID)
}

func TestUserService_Vote_GetUserError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepoInterface(ctrl)
	mockLogger := zaptest.NewLogger(t).Sugar()
	userService := NewUserService(mockRepo, mockLogger)

	testVote := &models.Vote{UserID: 1, ProfileID: 2, Value: 1}

	// Set expectations
	mockRepo.EXPECT().GetUserByID(gomock.Any(), testVote.UserID).Return(nil, errors.New("db error"))

	voteID, err := userService.Vote(context.Background(), testVote)
	assert.Error(t, err)
	assert.Equal(t, apperrors.InsertionFailedErr.Code, err.(*apperrors.AppError).Code)
	assert.Equal(t, constants.EmptyString, voteID)
}

func TestUserService_RevokeVote_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepoInterface(ctrl)
	mockLogger := zaptest.NewLogger(t).Sugar()
	userService := NewUserService(mockRepo, mockLogger)

	userID := uint(1)
	profileID := uint(2)

	// Set expectations
	mockRepo.EXPECT().DeleteVote(gomock.Any(), userID, profileID).Return(nil)

	err := userService.RevokeVote(context.Background(), userID, profileID)
	assert.NoError(t, err)
}

func TestUserService_RevokeVote_DeleteVoteError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepoInterface(ctrl)
	mockLogger := zaptest.NewLogger(t).Sugar()
	userService := NewUserService(mockRepo, mockLogger)

	userID := uint(1)
	profileID := uint(2)

	// Set expectations
	mockRepo.EXPECT().DeleteVote(gomock.Any(), userID, profileID).Return(errors.New("db error"))

	err := userService.RevokeVote(context.Background(), userID, profileID)
	assert.Error(t, err)
}
