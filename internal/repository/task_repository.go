package repository

import (
	"context"
	"github.com/advanced-coder-com/go-timekeeper/internal/db"
	"github.com/advanced-coder-com/go-timekeeper/internal/gormquery"
	"github.com/advanced-coder-com/go-timekeeper/internal/model"
	"gorm.io/gorm"
)

type TaskRepository interface {
	Create(ctx context.Context, task *model.Task) error
	GetByID(ctx context.Context, filters []gormquery.FilterGroup) (*model.Task, error)
	GetFilteredTasks(
		ctx context.Context,
		filters []gormquery.FilterGroup,
		options *gormquery.QueryOptions,
	) ([]model.Task, error)
	Update(ctx context.Context, task *model.Task) error
	Delete(ctx context.Context, task *model.Task) error
}

type taskRepository struct {
	database *gorm.DB
}

func NewTaskRepository() TaskRepository {
	return &taskRepository{database: db.Get()}
}

func (taskRepo *taskRepository) Create(ctx context.Context, task *model.Task) error {
	return taskRepo.database.WithContext(ctx).Create(task).Error
}

func (taskRepo *taskRepository) GetByID(ctx context.Context, filters []gormquery.FilterGroup) (*model.Task, error) {
	// Fixme get ID in params
	//filters := []gormquery.FilterGroup{
	//	gormquery.NewFilterGroup(
	//		gormquery.NewFilter("id", "=", taskID),
	//		gormquery.NewFilter("user_id", "=", userID),
	//	),
	//}
	var task model.Task
	query := taskRepo.database.WithContext(ctx).Model(&model.Task{})
	query = gormquery.ApplyFilters(query, filters)
	if err := query.First(&task).Error; err != nil {
		return nil, err
	}
	return &task, nil
}

func (taskRepo *taskRepository) GetFilteredTasks(
	ctx context.Context,
	filters []gormquery.FilterGroup,
	options *gormquery.QueryOptions,
) ([]model.Task, error) {
	var tasks []model.Task

	query := taskRepo.database.WithContext(ctx).Model(&model.Task{})
	query = gormquery.ApplyFilters(query, filters)

	if options != nil {
		query = gormquery.ApplyQueryOptions(query, *options)
	}

	if err := query.Find(&tasks).Error; err != nil {
		return nil, err
	}
	return tasks, nil
}

func (taskRepo *taskRepository) Update(ctx context.Context, task *model.Task) error {
	return taskRepo.database.WithContext(ctx).Save(task).Error
}

func (taskRepo *taskRepository) Delete(ctx context.Context, task *model.Task) error {
	return taskRepo.database.WithContext(ctx).Where("id = ?", task.ID).Delete(&model.Task{}).Error
}
