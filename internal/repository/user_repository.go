package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/advanced-coder-com/go-timekeeper/internal/db"
	"github.com/advanced-coder-com/go-timekeeper/internal/model"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	GetByID(ctx context.Context, id string) (*model.User, error)
	Update(ctx context.Context, user *model.User) error
	Delete(ctx context.Context, user *model.User) error
}

type userRepository struct {
	db *gorm.DB
}

const userRepoErrorPrefix = "UserRepository"

func NewUserRepository() UserRepository {
	return &userRepository{db: db.Get()}
}

func (repository *userRepository) Create(ctx context.Context, user *model.User) error {
	err := repository.db.WithContext(ctx).Create(user).Error
	if err != nil {
		err = fmt.Errorf("%s create user failed: %w", userRepoErrorPrefix, err)
	}
	return err
}

func (repository *userRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	err := repository.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		err = fmt.Errorf("%s get user by email failed: %w", userRepoErrorPrefix, err)
		return nil, err
	}
	return &user, nil
}

func (repository *userRepository) GetByID(ctx context.Context, id string) (*model.User, error) {
	var user model.User
	err := repository.db.WithContext(ctx).First(&user, "id = ?", id).Error
	if err != nil {
		err = fmt.Errorf("%s get user by id failed: %w", userRepoErrorPrefix, err)
		return nil, err
	}
	return &user, nil
}

func (repository *userRepository) Update(ctx context.Context, user *model.User) error {
	err := repository.db.WithContext(ctx).Save(user).Error
	if err != nil {
		err = fmt.Errorf("%s update user failed: %w", userRepoErrorPrefix, err)
	}
	return err
}

func (repository *userRepository) Delete(ctx context.Context, user *model.User) error {
	result := repository.db.WithContext(ctx).Delete(user)
	if result.Error != nil {
		return fmt.Errorf("%s delete user failed: %w", userRepoErrorPrefix, result.Error)
	}
	if result.RowsAffected == 0 {
		return errors.New(
			fmt.Sprintf("%s delete user failed: user you try to delete does not exist", userRepoErrorPrefix),
		)
	}
	return nil
}
