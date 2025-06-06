package service

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"strings"
	"time"

	"github.com/advanced-coder-com/go-timekeeper/internal/gormquery"
	"github.com/advanced-coder-com/go-timekeeper/internal/model"
	"github.com/advanced-coder-com/go-timekeeper/internal/repository"
	"github.com/google/uuid"
)

var (
	ErrProjectNotFound      = errors.New("project not found")
	ErrProjectAlreadyExists = errors.New("project with this name already exists for this user")
	ErrProjectDeleteFailed  = errors.New("failed to delete project")
	ErrProjectCreateFailed  = errors.New("failed to create project")
	ErrProjectUpdateFailed  = errors.New("failed to update project")
	ErrProjectListFailed    = errors.New("failed to list projects")
)

type ProjectInput struct {
	Name string `json:"name"`
}

type ProjectService struct {
	projectRepo repository.ProjectRepository
}

func NewProjectService() *ProjectService {
	return &ProjectService{
		projectRepo: repository.NewProjectRepository(),
	}
}

func (projectService *ProjectService) Create(ctx context.Context, userID string, input ProjectInput) (
	*model.Project,
	error,
) {
	filters := []gormquery.FilterGroup{
		gormquery.NewFilterGroup(
			gormquery.NewFilter("user_id", "=", userID),
			gormquery.NewFilter("LOWER(name)", "=", strings.ToLower(input.Name)),
		),
	}
	options := gormquery.QueryOptions{}

	existing, err := projectService.projectRepo.GetFilteredProjects(ctx, filters, options)
	if err != nil {
		return nil, err
	}
	if len(existing) > 0 {
		return nil, ErrProjectAlreadyExists
	}

	project := &model.Project{
		ID:        uuid.New(),
		Name:      input.Name,
		UserID:    uuid.MustParse(userID),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = projectService.projectRepo.Create(ctx, project)
	return project, err
}

func (projectService *ProjectService) GetAllByUser(ctx context.Context, userID string) ([]model.Project, error) {
	filters := []gormquery.FilterGroup{
		gormquery.NewFilterGroup(
			gormquery.NewFilter("user_id", "=", userID),
		),
	}
	return projectService.projectRepo.GetFilteredProjects(ctx, filters, gormquery.QueryOptions{})
}

func (projectService *ProjectService) GetByID(ctx context.Context, id string, userID string) (*model.Project, error) {
	filters := []gormquery.FilterGroup{
		gormquery.NewFilterGroup(
			gormquery.NewFilter("id", "=", id),
			gormquery.NewFilter("user_id", "=", userID),
		),
	}
	projects, err := projectService.projectRepo.GetFilteredProjects(ctx, filters, gormquery.QueryOptions{})
	if err != nil {
		return nil, err
	}
	if len(projects) == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return &projects[0], nil
}

func (projectService *ProjectService) Rename(ctx context.Context, id string, userID string, newName string) error {
	_, err := projectService.GetByID(ctx, id, userID)
	if err != nil {
		return err
	}

	updates := map[string]interface{}{"name": newName}
	return projectService.projectRepo.Update(ctx, id, updates)
}

func (projectService *ProjectService) Delete(ctx context.Context, projectID string) error {
	return projectService.projectRepo.DeleteByID(ctx, projectID)
}
