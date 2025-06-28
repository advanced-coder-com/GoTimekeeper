package repository

import (
	"context"
	"fmt"
	"github.com/advanced-coder-com/go-timekeeper/internal/db"
	"github.com/advanced-coder-com/go-timekeeper/internal/model"
	"gitlab.com/tozd/go/errors"
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

func NewUserRepository() UserRepository {
	return &userRepository{db: db.Get()}
}

func (repository *userRepository) Create(ctx context.Context, user *model.User) error {
	err := repository.db.WithContext(ctx).Create(user).Error
	if err != nil {
		err = errors.Errorf("create user failed: %v", err)
	}
	return err
}

func (repository *userRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	err := repository.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		err = errors.Errorf("get user by email failed: %v", err)
		return nil, err
	}
	return &user, nil
}

func (repository *userRepository) GetByID(ctx context.Context, id string) (*model.User, error) {
	var user model.User
	err := repository.db.WithContext(ctx).First(&user, "id = ?", id).Error
	if err != nil {
		err = errors.Errorf("get user by id failed: %v", err)
		return nil, err
	}
	return &user, nil
}

func (repository *userRepository) Update(ctx context.Context, user *model.User) error {
	err := repository.db.WithContext(ctx).Save(user).Error
	if err != nil {
		err = errors.Errorf("update user failed: %v", err)
	}
	return err
}

func (repository *userRepository) Delete(ctx context.Context, user *model.User) error {
	result := repository.db.WithContext(ctx).Delete(user)
	if result.Error != nil {
		return errors.Errorf("delete user failed: %v", result.Error)
	}
	if result.RowsAffected == 0 {
		return errors.New(
			fmt.Sprintf("delete user failed: user you try to delete does not exist. User ID: %s", user.ID),
		)
	}
	return nil
}
