package service

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"strings"
	"time"

	"github.com/advanced-coder-com/go-timekeeper/internal/model"
	"github.com/advanced-coder-com/go-timekeeper/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

// User related errors
var (
	ErrUserInvalidInput              = errors.New("email and password must be provided")
	ErrUserUnauthorized              = errors.New("unauthorized")
	ErrUserInvalidCredentials        = errors.New("invalid email or password")
	ErrUserEmailTaken                = errors.New("email already registered")
	ErrUserTokenInvalid              = errors.New("invalid or expired token")
	ErrUserMissingAuthHeader         = errors.New("missing authorization header")
	ErrUserInvalidAuthHeader         = errors.New("invalid authorization header")
	ErrUserMissingJWTSecret          = errors.New("missing JWT_SECRET")
	ErrUserTokenGenerationFailed     = errors.New("token generation failed")
	ErrUserChangePasswordInputFailed = errors.New("both old and new passwords are required")
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

func NewUserService() *UserService {
	repo := repository.NewUserRepository()
	return &UserService{repo: repo}
}

func (s *UserService) Signup(ctx context.Context, input UserInput) (*model.User, error) {
	err := input.validateUserInput()
	if err != nil {
		return nil, err
	}
	existing, _ := s.repo.FindByEmail(ctx, input.Email)
	if existing != nil {
		return nil, ErrUserEmailTaken
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
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
		return nil, err
	}
	user, err := s.repo.FindByEmail(ctx, input.Email)
	if err != nil {
		return nil, ErrUserInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		return nil, ErrUserInvalidCredentials
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
		return ErrUserChangePasswordInputFailed
	}

	user, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		return err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.OldPassword)); err != nil {
		return ErrUserInvalidCredentials
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
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
		return ErrUserInvalidInput
	}
	return nil
}
