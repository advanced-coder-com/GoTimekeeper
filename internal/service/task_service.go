package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/advanced-coder-com/go-timekeeper/internal/gormquery"
	"github.com/advanced-coder-com/go-timekeeper/internal/model"
	"github.com/advanced-coder-com/go-timekeeper/internal/repository"
	"github.com/google/uuid"
)

var (
	ErrTaskNotFound           = errors.New("task not found")
	ErrTaskAlreadyExists      = errors.New("task with this name already exists for this user")
	ErrTaskDeleteFailed       = errors.New("failed to delete task")
	ErrTaskCreateFailed       = errors.New("failed to create task")
	ErrTaskUpdateFailed       = errors.New("failed to update task")
	ErrTaskStartFailed        = errors.New("failed to start task")
	ErrTaskStopFailed         = errors.New("failed to stop task")
	ErrTaskListFailed         = errors.New("failed to list task")
	ErrTaskInvalidInput       = errors.New("invalid input")
	ErrTaskInvalidInputStatus = errors.New("invalid input status")
)

type TaskService struct {
	repo repository.TaskRepository
}

type CreateTaskInput struct {
	Name      string    `json:"name" binding:"required"`
	ProjectID uuid.UUID `json:"project_id,omitempty"`
	Tags      []string  `json:"tags"`
	Status    string    `json:"status"`
}

type UpdateTaskInput struct {
	Name      *string    `json:"name"`
	ProjectID *uuid.UUID `json:"project_id"`
	Tags      *[]string  `json:"tags"`
	Status    *string    `json:"status"`
}

func NewTaskService() *TaskService {
	return &TaskService{repo: repository.NewTaskRepository()}
}

func (taskService *TaskService) Create(ctx context.Context, userID string, input CreateTaskInput) (*model.Task, error) {
	var status model.TaskStatus
	if input.Status != "" {
		status = model.TaskStatus(input.Status)
	} else {
		status = model.DefaultStatus
	}
	task := &model.Task{
		UserID:    uuid.MustParse(userID),
		ProjectID: input.ProjectID,
		Name:      input.Name,
		Tags:      input.Tags,
		Status:    status,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	exists, err := taskService.checkExisting(ctx, task.UserID, task.ProjectID, task.Name)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrTaskAlreadyExists
	}
	err = taskService.repo.Create(ctx, task)
	return task, err
}

func (taskService *TaskService) GetAllByUser(ctx context.Context, userID string) ([]model.Task, error) {
	filters := []gormquery.FilterGroup{
		gormquery.NewFilterGroup(
			gormquery.NewFilter("user_id", "=", userID),
		),
	}
	return taskService.repo.GetFilteredTasks(ctx, filters, nil)
}

func (taskService *TaskService) GetAllActiveByUser(ctx context.Context, userID string) ([]model.Task, error) {
	filters := []gormquery.FilterGroup{
		gormquery.NewFilterGroup(
			gormquery.NewFilter("user_id", "=", userID),
			gormquery.NewFilter(
				"status",
				"IN",
				[]model.TaskStatus{model.StatusOpened, model.StatusWorkingOn},
			),
		),
	}
	return taskService.repo.GetFilteredTasks(ctx, filters, nil)
}

func (taskService *TaskService) GetByID(ctx context.Context, taskID int64, userID string) (*model.Task, error) {
	filters := []gormquery.FilterGroup{
		gormquery.NewFilterGroup(
			gormquery.NewFilter("id", "=", taskID),
			gormquery.NewFilter("user_id", "=", userID),
		),
	}
	return taskService.repo.GetByID(ctx, filters)
}

func (taskService *TaskService) Update(
	ctx context.Context,
	taskID int64,
	userID string,
	input UpdateTaskInput,
) (*model.Task, error) {
	task, err := taskService.GetByID(ctx, taskID, userID)
	if err != nil {
		return nil, err
	}

	if input.Name != nil && *input.Name != task.Name {
		task.Name = *input.Name
	}

	if input.ProjectID != nil && *input.ProjectID != task.ProjectID {
		task.ProjectID = *input.ProjectID
	}

	if input.Tags != nil && !equalStringSlices(*input.Tags, task.Tags) {
		task.Tags = *input.Tags
	}

	if input.Status != nil && string(task.Status) != *input.Status {
		if !model.IsValidTaskStatus(*input.Status) {
			return nil, fmt.Errorf("invalid status")
		}
		task.Status = model.TaskStatus(*input.Status)
	}

	task.UpdatedAt = time.Now()

	exists, err := taskService.checkExisting(ctx, task.UserID, task.ProjectID, task.Name)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrTaskAlreadyExists
	}
	err = taskService.repo.Update(ctx, task)

	return task, err
}

func (taskService *TaskService) Delete(ctx context.Context, id int64, userID string) error {
	task, err := taskService.GetByID(ctx, id, userID)
	if err != nil {
		return err
	}
	return taskService.repo.Delete(ctx, task)
}

func (taskService *TaskService) Start(ctx context.Context, id int64, userID string) error {
	task, err := taskService.GetByID(ctx, id, userID)
	if err != nil {
		return err
	}
	err = checkIfTaskIsNotClosed(task)
	if err != nil {
		return err
	}
	task.Status = model.StatusWorkingOn
	task.UpdatedAt = time.Now()
	// TODO: Create timeslot after Timeslot entity is ready
	return taskService.repo.Update(ctx, task)
}

func (taskService *TaskService) Stop(ctx context.Context, id int64, userID string) error {
	task, err := taskService.GetByID(ctx, id, userID)
	if err != nil {
		return err
	}
	err = checkIfTaskIsNotClosed(task)
	if err != nil {
		return err
	}
	task.Status = model.StatusClosed
	task.UpdatedAt = time.Now()
	// TODO: Stop timeslot after Timeslot entity is ready
	return taskService.repo.Update(ctx, task)
}

func (taskService *TaskService) StopAll(ctx context.Context, userID string) error {
	tasks, err := taskService.GetAllActiveByUser(ctx, userID)
	if err != nil {
		return err
	}
	for _, task := range tasks {
		err = checkIfTaskIsNotClosed(&task)
		if err != nil {
			return err
		}
		task.Status = model.StatusOpened
		task.UpdatedAt = time.Now()
		// TODO: Stop timeslot after Timeslot entity is ready
		err := taskService.repo.Update(ctx, &task)
		if err != nil {
			return err
		}
	}
	return nil
}

func (taskService *TaskService) Close(ctx context.Context, id int64, userID string) error {
	task, err := taskService.GetByID(ctx, id, userID)
	if err != nil {
		return err
	}
	err = checkIfTaskIsNotClosed(task)
	if err != nil {
		return err
	}
	task.Status = model.StatusClosed
	task.UpdatedAt = time.Now()
	// TODO: Stop timeslot after Timeslot entity is ready
	return taskService.repo.Update(ctx, task)
}

func (taskService *TaskService) checkExisting(
	ctx context.Context,
	userID uuid.UUID,
	projectID uuid.UUID,
	name string,
) (bool, error) {
	filters := []gormquery.FilterGroup{
		gormquery.NewFilterGroup(
			gormquery.NewFilter("user_id", "=", userID),
			gormquery.NewFilter("project_id", "=", projectID),

			gormquery.NewFilter("LOWER(name)", "=", strings.ToLower(name)),
		),
	}

	tasks, err := taskService.repo.GetFilteredTasks(ctx, filters, nil)
	if err != nil {
		return false, err
	}
	return len(tasks) == 0, nil
}

func checkIfTaskIsNotClosed(task *model.Task) error {
	if task.Status == model.StatusClosed {
		return errors.New("requested task is closed")
	}
	return nil
}

func equalStringSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
