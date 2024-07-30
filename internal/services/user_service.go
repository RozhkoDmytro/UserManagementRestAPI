package services

import (
	"context"
	"strconv"

	"gitlab.com/jkozhemiaka/web-layout/internal/constants"

	"gitlab.com/jkozhemiaka/web-layout/internal/models"
	"gitlab.com/jkozhemiaka/web-layout/internal/repositories"
	"go.uber.org/zap"
)

type UserService struct {
	userRepo *repositories.UserRepo
	logger   *zap.SugaredLogger
}

type UserServiceInterface interface {
	CreateUser(ctx context.Context, user *models.User) (string, error)
	DeleteUser(ctx context.Context, userID string) (*models.User, error)
	GetUser(ctx context.Context, userID string) (*models.User, error)
	UpdateUser(ctx context.Context, userID string, user *models.User) (*models.User, error)
	ListUsers(ctx context.Context, page, pageSize string) ([]models.User, error)
}

func NewUserService(userRepo *repositories.UserRepo, logger *zap.SugaredLogger) UserServiceInterface {
	return &UserService{
		userRepo: userRepo,
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

func (service *UserService) ListUsers(ctx context.Context, page, pageSize string) (user []models.User, err error) {
	intPage, err := strconv.Atoi(page)
	if err != nil {
		service.logger.Error(err)
		return nil, err
	}

	intPageSize, err := strconv.Atoi(pageSize)
	if err != nil {
		service.logger.Error(err)
		return nil, err
	}

	user, err = service.userRepo.ListUsers(ctx, intPage, intPageSize)
	if err != nil {
		service.logger.Error(err)
		return nil, err
	}

	return user, nil
}
