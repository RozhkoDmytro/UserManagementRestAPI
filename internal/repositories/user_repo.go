package repositories

import (
	"context"
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
	result := tx.First(&user, "id = ? AND deleted_at = ?", userID, time.Time{})
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
	result := tx.First(&user, "id = ? AND deleted_at = ?", userID, time.Time{})
	if result.Error != nil {
		if result.RowsAffected == 0 {
			repo.logger.Warn("No user found with the given ID.")
			return nil, apperrors.NoRecordFoundErr.AppendMessage("No user found with the given ID.")
		}
		repo.logger.Error(result.Error)
		return nil, apperrors.DeletionFailedErr.AppendMessage(result.Error.Error())
	}

	// Apply the updates
	if updatedData.FirstName != "" {
		user.FirstName = updatedData.FirstName
	}
	if updatedData.LastName != "" {
		user.LastName = updatedData.LastName
	}
	if updatedData.Password != "" {
		user.Password = updatedData.Password
	}
	if updatedData.Email != "" {
		user.Email = updatedData.Email
	}
	if !updatedData.DeletedAt.IsZero() {
		user.DeletedAt = updatedData.DeletedAt
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

	result := tx.Limit(pageSize).Offset(offset).Find(&users, "deleted_at = ?", time.Time{})
	if result.Error != nil {
		repo.logger.Error(result.Error)
		return nil, apperrors.DeletionFailedErr.AppendMessage(result.Error.Error())
	}

	return users, nil
}

func (repo *UserRepo) CountUsers(ctx context.Context) (int, error) {
	var count int64
	tx := repo.db.WithContext(ctx)
	result := tx.Model(&models.User{}).Where("deleted_at = ?", time.Time{}).Count(&count)
	if result.Error != nil {
		repo.logger.Error(result.Error)
		return 0, apperrors.DeletionFailedErr.AppendMessage(result.Error.Error())
	}
	return int(count), nil
}
