package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/advanced-coder-com/go-timekeeper/internal/db"

	"github.com/advanced-coder-com/go-timekeeper/internal/model"
	"gorm.io/gorm"
)

const userRepositoryErrorPrefix = "User Repository Error"

type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	FindByEmail(ctx context.Context, email string) (*model.User, error)
	DeleteByID(ctx context.Context, id string) error
	FindByID(ctx context.Context, id string) (*model.User, error)
	Update(ctx context.Context, user *model.User) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository() UserRepository {
	return &userRepository{db: db.Get()}
}

func (repository *userRepository) Create(ctx context.Context, user *model.User) error {
	return repository.db.WithContext(ctx).Create(user).Error
}

func (repository *userRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	err := repository.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (repository *userRepository) FindByID(ctx context.Context, id string) (*model.User, error) {
	var user model.User
	err := repository.db.WithContext(ctx).First(&user, "id = ?", id).Error
	if err != nil {
		return nil, wrapError(&err)
	}
	return &user, nil
}

func (repository *userRepository) Update(ctx context.Context, user *model.User) error {
	return repository.db.WithContext(ctx).Save(user).Error
}

func (repository *userRepository) DeleteByID(ctx context.Context, id string) error {

	result := repository.db.WithContext(ctx).Where("id = ?", id).Delete(&model.User{})
	if result.Error != nil {
		return wrapError(&result.Error)
	}
	if result.RowsAffected == 0 {
		return errors.New("user not found")
	}
	return nil
}

// Wrap Database error with the prefix
func wrapError(err *error) error {
	return fmt.Errorf("%s: %w", userRepositoryErrorPrefix, *err)
}
