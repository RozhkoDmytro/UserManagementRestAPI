package repositories

import (
	"context"
	"errors"
	"time"

	"gitlab.com/jkozhemiaka/web-layout/internal/apperrors"

	"gitlab.com/jkozhemiaka/web-layout/internal/models"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type UserRepo struct {
	db     *gorm.DB
	logger *zap.SugaredLogger
}

type UserRepoInterface interface {
	CreateUser(ctx context.Context, user *models.User) (*models.User, error)
	GetUser(ctx context.Context, userID string) (*models.User, error)
	DeleteUser(ctx context.Context, userID string) (*models.User, error)
	UpdateUser(ctx context.Context, userID string, updatedData *models.User) (*models.User, error)
	ListUsers(ctx context.Context, page int, pageSize int) ([]models.User, error)
	CountUsers(ctx context.Context) (int, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	GetUserByID(ctx context.Context, userID uint, user *models.User) error
	GetVote(ctx context.Context, userID uint, profileID uint) (*models.Vote, error)
	CreateVote(ctx context.Context, vote *models.Vote) (*models.Vote, error)
	UpdateVote(ctx context.Context, vote *models.Vote) error
	DeleteVote(ctx context.Context, userID uint, profileID uint) error
}

func NewUserRepo(db *gorm.DB, logger *zap.SugaredLogger) *UserRepo {
	return &UserRepo{
		db:     db,
		logger: logger,
	}
}

func (repo *UserRepo) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	tx := repo.db.WithContext(ctx)
	tx.Create(user)
	if tx.Error != nil {
		repo.logger.Error(tx.Error)
		return nil, apperrors.InsertionFailedErr.AppendMessage(tx.Error)
	}

	return user, nil
}

func (repo *UserRepo) GetUser(ctx context.Context, userID string) (*models.User, error) {
	tx := repo.db.WithContext(ctx)
	var user models.User

	// Fetch the user to be updated
	result := tx.First(&user, "id = ? AND (deleted_at IS NULL OR deleted_at = ?)", userID, time.Time{})
	if result.Error != nil {
		if result.RowsAffected == 0 {
			repo.logger.Warn("No user found with the given ID.")
			return nil, apperrors.NoRecordFoundErr.AppendMessage("No user found with the given ID.")
		}
		repo.logger.Error(result.Error)
		return nil, apperrors.DeletionFailedErr.AppendMessage(result.Error.Error())
	}

	return &user, nil
}

func (repo *UserRepo) DeleteUser(ctx context.Context, userID string) (*models.User, error) {
	return repo.UpdateUser(ctx, userID, &models.User{DeletedAt: time.Now()})
}

func (repo *UserRepo) UpdateUser(ctx context.Context, userID string, updatedData *models.User) (*models.User, error) {
	tx := repo.db.WithContext(ctx)
	var user models.User

	// Fetch the user to be updated
	result := tx.First(&user, "id = ? AND (deleted_at IS NULL OR deleted_at = ?)", userID, time.Time{})
	if result.Error != nil {
		if result.RowsAffected == 0 {
			repo.logger.Warn("No user found with the given ID.")
			return nil, apperrors.NoRecordFoundErr.AppendMessage("No user found with the given ID.")
		}
		repo.logger.Error(result.Error)
		return nil, apperrors.DeletionFailedErr.AppendMessage(result.Error.Error())
	}

	// Apply the updates
	// Check email uniqueness if it changes
	if updatedData.Email != "" && updatedData.Email != user.Email {
		var existingUser models.User
		result := tx.First(&existingUser, "email = ?", updatedData.Email)
		if result.RowsAffected > 0 {
			repo.logger.Warn("The email is already occupied by another user.")
			return nil, apperrors.DeletionFailedErr.AppendMessage("The email is already occupied by another user.")
		}
		user.Email = updatedData.Email
	}

	if updatedData.FirstName != "" {
		user.FirstName = updatedData.FirstName
	}
	if updatedData.LastName != "" {
		user.LastName = updatedData.LastName
	}
	if updatedData.Password != "" {
		user.Password = updatedData.Password
	}
	if !updatedData.DeletedAt.IsZero() {
		user.DeletedAt = updatedData.DeletedAt
	}
	if updatedData.RoleID > 0 {
		user.RoleID = updatedData.RoleID
	}
	// Add other fields to update as needed

	// Save the changes
	result = tx.Save(&user)
	if result.Error != nil {
		repo.logger.Error(result.Error)
		return nil, apperrors.DeletionFailedErr.AppendMessage(result.Error.Error())
	}

	return &user, nil
}

func (repo *UserRepo) ListUsers(ctx context.Context, page int, pageSize int) ([]models.User, error) {
	var users []models.User
	tx := repo.db.WithContext(ctx)

	// Calculate offset for pagination
	offset := (page - 1) * pageSize

	result := tx.Limit(pageSize).Offset(offset).Preload("Role").Find(&users, "deleted_at IS NULL OR deleted_at = ?", time.Time{})
	if result.Error != nil {
		repo.logger.Error(result.Error)
		return nil, apperrors.DeletionFailedErr.AppendMessage(result.Error.Error())
	}

	return users, nil
}

func (repo *UserRepo) CountUsers(ctx context.Context) (int, error) {
	var count int64
	tx := repo.db.WithContext(ctx)
	result := tx.Model(&models.User{}).Where("deleted_at IS NULL OR deleted_at = ?", time.Time{}).Count(&count)
	if result.Error != nil {
		repo.logger.Error(result.Error)
		return 0, apperrors.DeletionFailedErr.AppendMessage(result.Error.Error())
	}
	return int(count), nil
}

func (repo *UserRepo) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	tx := repo.db.WithContext(ctx).
		Where("email = ? AND (deleted_at IS NULL OR deleted_at = ?)", email, time.Time{}).
		Preload("Role").
		First(&user)
	if tx.Error != nil {
		if tx.RowsAffected == 0 {
			return nil, nil // No user found
		}
		repo.logger.Error(tx.Error)
		return nil, tx.Error
	}
	return &user, nil
}

func (repo *UserRepo) GetUserByID(ctx context.Context, userID uint, user *models.User) error {
	result := repo.db.WithContext(ctx).First(user, userID)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return apperrors.NoRecordFoundErr.AppendMessage("User not found.")
		}
		return result.Error
	}
	return nil
}

func (repo *UserRepo) GetVote(ctx context.Context, userID uint, profileID uint) (*models.Vote, error) {
	var vote models.Vote
	result := repo.db.WithContext(ctx).Where("user_id = ? AND profile_id = ?", userID, profileID).First(&vote)
	if result.Error != nil {
		return nil, result.Error
	}
	return &vote, nil
}

func (repo *UserRepo) CreateVote(ctx context.Context, vote *models.Vote) (*models.Vote, error) {
	if err := repo.db.WithContext(ctx).Create(vote).Error; err != nil {
		repo.logger.Error("Failed to create vote", zap.Error(err))
		return nil, err
	}
	return vote, nil
}

func (repo *UserRepo) UpdateVote(ctx context.Context, vote *models.Vote) error {
	if err := repo.db.WithContext(ctx).Save(vote).Error; err != nil {
		repo.logger.Error("Failed to update vote", zap.Error(err))
		return err
	}
	return nil
}

func (repo *UserRepo) DeleteVote(ctx context.Context, userID uint, profileID uint) error {
	tx := repo.db.WithContext(ctx)

	// Find the vote
	var vote models.Vote
	result := tx.Where("user_id = ? AND profile_id = ?", userID, profileID).First(&vote)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return apperrors.NoRecordFoundErr.AppendMessage("Vote not found.")
		}
		return result.Error
	}

	// Delete the vote
	if err := tx.Delete(&vote).Error; err != nil {
		return apperrors.DeletionFailedErr.AppendMessage(err.Error())
	}

	return nil
}
