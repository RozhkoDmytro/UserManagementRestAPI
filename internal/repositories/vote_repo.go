package repositories

import (
	"context"
	"errors"

	"gitlab.com/jkozhemiaka/web-layout/internal/apperrors"
	"gitlab.com/jkozhemiaka/web-layout/internal/models"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

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

func (repo *UserRepo) UpdateVote(ctx context.Context, vote *models.Vote) (*models.Vote, error) {
	if err := repo.db.WithContext(ctx).Save(vote).Error; err != nil {
		repo.logger.Error("Failed to update vote", zap.Error(err))
		return nil, err
	}
	return vote, nil
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
