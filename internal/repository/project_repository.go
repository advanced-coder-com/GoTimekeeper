package repository

import (
	"context"

	"github.com/advanced-coder-com/go-timekeeper/internal/db"
	"github.com/advanced-coder-com/go-timekeeper/internal/gormquery"
	"github.com/advanced-coder-com/go-timekeeper/internal/model"
	"gorm.io/gorm"
)

type ProjectRepository interface {
	Create(ctx context.Context, project *model.Project) error
	Update(ctx context.Context, id string, updates map[string]interface{}) error
	// FIXME options should be *gormquery.QueryOptions
	GetFilteredProjects(ctx context.Context, filters []gormquery.FilterGroup, options gormquery.QueryOptions) (
		[]model.Project,
		error,
	)
	DeleteByID(ctx context.Context, id string) error
}

type projectRepository struct {
	database *gorm.DB
}

func NewProjectRepository() ProjectRepository {
	return &projectRepository{database: db.Get()}
}

func (projectRepo *projectRepository) Create(ctx context.Context, project *model.Project) error {
	return projectRepo.database.WithContext(ctx).Create(project).Error
}

func (projectRepo *projectRepository) GetFilteredProjects(
	ctx context.Context,
	filters []gormquery.FilterGroup,
	options gormquery.QueryOptions,
) ([]model.Project, error) {
	var projects []model.Project

	query := projectRepo.database.WithContext(ctx).Model(&model.Project{})
	query = gormquery.ApplyFilters(query, filters)
	query = gormquery.ApplyQueryOptions(query, options)

	err := query.Find(&projects).Error
	return projects, err
}

func (projectRepo *projectRepository) Update(ctx context.Context, id string, updates map[string]interface{}) error {
	// FIXME param should be  *model.Project
	result := projectRepo.database.WithContext(ctx).
		Model(&model.Project{}).
		Where("id = ?", id).
		Updates(updates)

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (projectRepo *projectRepository) DeleteByID(ctx context.Context, id string) error {
	result := projectRepo.database.WithContext(ctx).Where("id = ?", id).Delete(&model.Project{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
