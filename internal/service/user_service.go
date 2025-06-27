package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"strings"
	"time"

	"github.com/advanced-coder-com/go-timekeeper/internal/model"
	"github.com/advanced-coder-com/go-timekeeper/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

// User related errors
var (
	ErrUserUnauthorized         = errors.New("user is unauthorized")
	ErrUserTokenInvalid         = errors.New("invalid or expired token")
	ErrUserMissingAuthHeader    = errors.New("missing authorization header")
	ErrUserInvalidAuthHeader    = errors.New("invalid authorization header")
	ErrUserMissingJWTSecret     = errors.New("missing JWT_SECRET")
	ErrUserInvalidInput         = errors.New("invalid input")
	ErrUserSignInFailed         = errors.New("user sign in failed")
	ErrUserSignUpFailed         = errors.New("user sign up failed")
	ErrGetUserFailed            = errors.New("cannot get user with provided credentials")
	ErrUserDeleteFailed         = errors.New("cannot delete user")
	ErrUserChangePasswordFailed = errors.New("changing password failed")
)

// UserInput Input for user API routes
type UserInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// ChangePasswordInput Input for changing password
type ChangePasswordInput struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

type UserService struct {
	repo repository.UserRepository
}

const userServiceLogPrefix = "UserService"

func NewUserService() *UserService {
	repo := repository.NewUserRepository()
	return &UserService{repo: repo}
}

func (s *UserService) Signup(ctx context.Context, input UserInput) (*model.User, error) {
	err := input.validateUserInput()
	if err != nil {
		err = fmt.Errorf("%s: %w", userServiceLogPrefix, err)
		return nil, err
	}
	existing, _ := s.repo.FindByEmail(ctx, input.Email)
	if existing != nil {
		return nil, errors.New(fmt.Sprintf("%s user with this email exists", userServiceLogPrefix))
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", userServiceLogPrefix, err)
	}

	user := &model.User{
		ID:        uuid.New(),
		Email:     input.Email,
		Password:  string(hashedPassword),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) Signin(ctx context.Context, input UserInput) (*model.User, error) {
	err := input.validateUserInput()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", userServiceLogPrefix, err)

	}
	user, err := s.repo.FindByEmail(ctx, input.Email)
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		return nil, fmt.Errorf("%s: %w", userServiceLogPrefix, err)
	}
	return user, nil
}

func (s *UserService) GetUser(ctx context.Context, userId string) (*model.User, error) {
	user, err := s.repo.FindByID(ctx, userId)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) ChangePassword(ctx context.Context, userID string, input ChangePasswordInput) error {
	if input.OldPassword == "" || input.NewPassword == "" {
		return errors.New(fmt.Sprintf("%s both old and new passwords are required", userServiceLogPrefix))
	}
	if input.OldPassword == input.NewPassword {
		return errors.New(fmt.Sprintf("%s Old password must not be same as a new one", userServiceLogPrefix))
	}

	user, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		return err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.OldPassword)); err != nil {
		return fmt.Errorf("%s: %w", userServiceLogPrefix, err)
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("%s: %w", userServiceLogPrefix, err)
	}

	user.Password = string(hashed)
	user.UpdatedAt = time.Now()

	return s.repo.Update(ctx, user)
}

func (s *UserService) Delete(ctx context.Context, userId string) error {
	return s.repo.DeleteByID(ctx, userId)
}

// validation of user input
func (input *UserInput) validateUserInput() error {
	input.Email = strings.TrimSpace(strings.ToLower(input.Email))
	if input.Email == "" || input.Password == "" {
		return errors.New("email and password must be provided")
	}
	return nil
}
