package service

import (
	"context"
	"fmt"
	"github.com/advanced-coder-com/go-timekeeper/internal/validator"
	"github.com/google/uuid"
	"gitlab.com/tozd/go/errors"
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

func (userService *UserService) Signup(ctx context.Context, input UserInput) (*model.User, error) {
	err := input.validateUserInput()
	if err != nil {
		return nil, WrapPublicMessage(err, err.Error())
	}
	existing, _ := userService.repo.GetByEmail(ctx, input.Email)
	if existing != nil {
		message := fmt.Sprintf("User with email %s already exists", input.Email)

		err = errors.New(message)
		return nil, WrapPublicMessage(err, message)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.Errorf("%v", err)
	}

	user := &model.User{
		ID:        uuid.New(),
		Email:     input.Email,
		Password:  string(hashedPassword),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := userService.repo.Create(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

func (userService *UserService) Signin(ctx context.Context, input UserInput) (*model.User, error) {
	err := input.validateUserInput()
	if err != nil {
		return nil, err

	}
	user, err := userService.repo.GetByEmail(ctx, input.Email)
	if err != nil {
		return nil, WrapPublicMessage(err, "User with provided email does not exist")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		return nil, errors.Errorf("%v", err)
	}
	return user, nil
}

func (userService *UserService) GetUser(ctx context.Context, userId string) (*model.User, error) {
	user, err := userService.repo.GetByID(ctx, userId)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (userService *UserService) ChangePassword(ctx context.Context, userID string, input ChangePasswordInput) error {
	if input.OldPassword == "" || input.NewPassword == "" {
		err := errors.New("both old and new passwords are required")
		return WrapPublicMessage(err, err.Error())
	}
	if input.OldPassword == input.NewPassword {
		err := errors.New("Old password must not be same as a new one")
		return WrapPublicMessage(err, err.Error())
	}

	user, err := userService.repo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.OldPassword)); err != nil {
		return errors.Errorf("%v", err)
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return errors.Errorf("%v", err)
	}

	user.Password = string(hashed)
	user.UpdatedAt = time.Now()

	return userService.repo.Update(ctx, user)
}

func (userService *UserService) Delete(ctx context.Context, userId string) error {
	user, err := userService.repo.GetByID(ctx, userId)
	if err != nil {
		return err
	}
	err = userService.repo.Delete(ctx, user)
	if err != nil {
		return err
	}
	return nil
}

// validation of user input
func (input *UserInput) validateUserInput() error {
	err := validator.ValidateEmail(input.Email)
	if err != nil {
		return errors.New("invalid email")
	}
	err = validator.ValidatePassword(input.Password)
	if err != nil {
		return errors.Errorf("%v", err)
	}
	return nil
}
