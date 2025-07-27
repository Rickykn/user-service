package service

import (
	"context"
	"github.com/Rickykn/user-service/src/model"
	"github.com/Rickykn/user-service/src/repository"
	"github.com/Rickykn/user-service/src/utils"
	"time"
)

type IUserService interface {
	GetUserByUsername(ctx context.Context, username string) (*model.User, error)
	CreateUser(ctx context.Context, user *model.User) error
}

type userService struct {
	repo repository.IUserRepository
}

func NewUserService(repo repository.IUserRepository) IUserService {
	return &userService{repo: repo}
}

func (u *userService) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	byUsername, err := u.repo.GetByUsername(ctx, username)

	if err != nil {
		panic("error get database")
	}
	return byUsername, nil
}

func (u *userService) CreateUser(ctx context.Context, user *model.User) error {
	password, err := utils.HashPassword(user.Password)

	if err != nil {
		panic("error hasing password")
	}

	newUser := &model.User{
		Username:  user.Username,
		Email:     user.Email,
		Password:  password,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err = u.repo.CreateUser(ctx, newUser)
	if err != nil {
		panic("error create user from database")
	}

	return nil
}
