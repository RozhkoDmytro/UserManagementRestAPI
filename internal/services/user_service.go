package services

import (
	"context"
	"strconv"

	"gitlab.com/jkozhemiaka/web-layout/internal/constants"
	"gitlab.com/jkozhemiaka/web-layout/internal/repositories"

	"gitlab.com/jkozhemiaka/web-layout/internal/models"
	"go.uber.org/zap"
)

type UserService struct {
	userRepo repositories.UserRepoInterface
	logger   *zap.SugaredLogger
}

type UserServiceInterface interface {
	CreateUser(ctx context.Context, user *models.User) (string, error)
	DeleteUser(ctx context.Context, userID string) (*models.User, error)
	GetUser(ctx context.Context, userID string) (*models.User, error)
	UpdateUser(ctx context.Context, userID string, user *models.User) (*models.User, error)
	ListUsers(ctx context.Context, page, pageSize int) ([]models.User, error)
	CountUsers(ctx context.Context) (int, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
}

func NewUserService(userRepo repositories.UserRepoInterface, logger *zap.SugaredLogger) UserServiceInterface {
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
