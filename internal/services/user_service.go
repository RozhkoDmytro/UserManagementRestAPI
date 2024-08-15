package services

import (
	"context"
	"time"

	"gitlab.com/jkozhemiaka/web-layout/internal/apperrors"
	"gitlab.com/jkozhemiaka/web-layout/internal/repositories"
	"gorm.io/gorm"

	"gitlab.com/jkozhemiaka/web-layout/internal/models"
	"go.uber.org/zap"
)

type UserService struct {
	userRepo repositories.UserRepoInterface
	logger   *zap.SugaredLogger
}

type UserServiceInterface interface {
	CreateUser(ctx context.Context, user *models.User) (uint, error)
	DeleteUser(ctx context.Context, userID string) (*models.User, error)
	GetUser(ctx context.Context, userID string) (*models.User, error)
	UpdateUser(ctx context.Context, userID string, user *models.User) (*models.User, error)
	ListUsers(ctx context.Context, page, pageSize int) ([]models.User, error)
	CountUsers(ctx context.Context) (int, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	Vote(ctx context.Context, vote *models.Vote) (uint, error)
	RevokeVote(ctx context.Context, userID uint, profileID uint) error
}

func NewUserService(userRepo repositories.UserRepoInterface, logger *zap.SugaredLogger) UserServiceInterface {
	return &UserService{
		userRepo: userRepo,
		logger:   logger,
	}
}

func (service *UserService) CreateUser(ctx context.Context, user *models.User) (userId uint, err error) {
	insertedUser, err := service.userRepo.CreateUser(ctx, user)
	if err != nil {
		service.logger.Error(err)
		return 0, err
	}

	return insertedUser.ID, nil
}

func (service *UserService) GetUser(ctx context.Context, userID string) (user *models.User, err error) {
	user, err = service.userRepo.GetUser(ctx, userID)
	if err != nil {
		service.logger.Error(err)
		return nil, err
	}

	return user, nil
}

func (service *UserService) DeleteUser(ctx context.Context, userID string) (user *models.User, err error) {
	user, err = service.userRepo.DeleteUser(ctx, userID)
	if err != nil {
		service.logger.Error(err)
		return nil, err
	}

	return user, nil
}

func (service *UserService) UpdateUser(ctx context.Context, userID string, updatedData *models.User) (user *models.User, err error) {
	user, err = service.userRepo.UpdateUser(ctx, userID, updatedData)
	if err != nil {
		service.logger.Error(err)
		return nil, err
	}

	return user, nil
}

func (service *UserService) ListUsers(ctx context.Context, page, pageSize int) (user []models.User, err error) {
	user, err = service.userRepo.ListUsers(ctx, page, pageSize)
	if err != nil {
		service.logger.Error(err)
		return nil, err
	}

	return user, nil
}

func (service *UserService) CountUsers(ctx context.Context) (int, error) {
	count, err := service.userRepo.CountUsers(ctx)
	if err != nil {
		service.logger.Error(err)
		return 0, err
	}

	return count, nil
}

func (service *UserService) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	user, err := service.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		service.logger.Error(err)
		return nil, err
	}

	return user, nil
}

func (service *UserService) Vote(ctx context.Context, vote *models.Vote) (uint, error) {
	// Get the user profile
	var user *models.User
	user, err := service.userRepo.GetUserByID(ctx, vote.UserID)
	if err != nil {
		service.logger.Error("Failed to get user", zap.Error(err))
		return 0, apperrors.InsertionFailedErr.AppendMessage(err.Error())
	}

	// Check if the user has voted within the last hour
	if time.Since(user.VoteUpdatedAt) < time.Hour {
		return 0, &apperrors.VoteCooldownErr
	}

	// Check if the user has already voted for this profile
	existingVote, err := service.userRepo.GetVote(ctx, vote.UserID, vote.ProfileID)
	if err != nil && err != gorm.ErrRecordNotFound {
		service.logger.Error("Failed to check existing vote", zap.Error(err))
		return 0, apperrors.InsertionFailedErr.AppendMessage(err.Error())
	}

	if existingVote != nil {
		// Update existing vote
		existingVote.Value = vote.Value
		_, err = service.userRepo.UpdateVote(ctx, existingVote)
		if err != nil {
			service.logger.Error("Failed to update vote", zap.Error(err))
			return 0, apperrors.UpdateFailedErr.AppendMessage(err.Error())
		}
		return existingVote.ID, nil
	}

	// Create new vote
	insertedVote, err := service.userRepo.CreateVote(ctx, vote)
	if err != nil {
		service.logger.Error("Failed to create vote", zap.Error(err))
		return 0, apperrors.InsertionFailedErr.AppendMessage(err.Error())
	}

	return insertedVote.ID, nil
}

func (service *UserService) RevokeVote(ctx context.Context, userID uint, profileID uint) error {
	// Proceed to delete the vote
	err := service.userRepo.DeleteVote(ctx, userID, profileID)
	if err != nil {
		return err
	}

	return nil
}
