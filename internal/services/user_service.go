package services

import (
	"context"
	"strconv"

	"gitlab.com/jkozhemiaka/web-layout/internal/constants"

	"gitlab.com/jkozhemiaka/web-layout/internal/models"
	"gitlab.com/jkozhemiaka/web-layout/internal/repositories"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type UserService struct {
	userRepo *repositories.UserRepo
	logger   *zap.SugaredLogger
}

func NewUserService(db *gorm.DB, logger *zap.SugaredLogger) *UserService {
	return &UserService{
		userRepo: repositories.NewUserRepo(db, logger),
		logger:   logger,
	}
}

func (service *UserService) CreateUser(ctx context.Context, user *models.User) (userId string, err error) {
	insertedUser, err := service.userRepo.CreateUser(ctx, user)
	if err != nil {
		service.logger.Error(err)
		return constants.EmptyString, err
	}

	return strconv.Itoa(int(insertedUser.ID)), nil
}

func (service *UserService) GetUser(ctx context.Context, userID string) (user *models.User, err error) {
	user, err = service.userRepo.GetUser(ctx, userID)
	if err != nil {
		service.logger.Error(err)
		return nil, err
	}

	return user, nil
}
