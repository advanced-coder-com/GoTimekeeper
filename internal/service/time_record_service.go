package service

import (
	"context"
	"errors"
	"github.com/advanced-coder-com/go-timekeeper/internal/gormquery"
	"github.com/advanced-coder-com/go-timekeeper/internal/model"
	"github.com/advanced-coder-com/go-timekeeper/internal/repository"
	"github.com/google/uuid"
	"time"
)

var (
	ErrTimeRecordCreationFailed = errors.New("cannot create time record fot task has not closed time record")
	ErrTimeRecordMultipleActive = errors.New("task has more that one active time record")
)

type TimeRecordService struct {
	repo repository.TimeRecordRepository
}

type UpdateTimeRecordInput struct {
	TaskID    *uint64
	StartTime *time.Time
	EndTime   *time.Time
	IsClosed  *bool
}

func NewTimeRecordService() *TimeRecordService {
	return &TimeRecordService{repo: repository.NewTimeRecordRepository()}
}

func (timeRecordService *TimeRecordService) Create(
	ctx context.Context,
	userID string,
	taskID uint64,
) (*model.TimeRecord, error) {

	timeRecord := &model.TimeRecord{
		UserID:    uuid.MustParse(userID),
		TaskID:    taskID,
		StartTime: time.Now(),
	}
	err := timeRecordService.createTimeRecordValidate(ctx, timeRecord)
	if err != nil {
		return nil, err
	}

	err = timeRecordService.repo.Create(ctx, timeRecord)
	return timeRecord, err
}

func (timeRecordService *TimeRecordService) GetByID(ctx context.Context, id uint64) (*model.TimeRecord, error) {
	return timeRecordService.repo.GetByID(ctx, id)
}

func (timeRecordService *TimeRecordService) GetByTaskID(ctx context.Context, taskID uint64) (*[]model.TimeRecord, error) {
	return timeRecordService.repo.GetByTaskID(ctx, taskID)
}

func (timeRecordService *TimeRecordService) Update(
	ctx context.Context,
	id uint64,
	input UpdateTimeRecordInput,
) (*model.TimeRecord, error) {
	timeRecord, err := timeRecordService.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if input.TaskID != nil && *input.TaskID != timeRecord.TaskID {
		timeRecord.TaskID = *input.TaskID
	}

	if input.StartTime != nil && *input.StartTime != timeRecord.StartTime {
		timeRecord.StartTime = *input.StartTime
	}

	if input.EndTime != nil && *input.EndTime != *timeRecord.EndTime {
		timeRecord.EndTime = input.EndTime
	}

	if input.IsClosed != nil && *input.IsClosed != timeRecord.IsClosed {
		timeRecord.IsClosed = *input.IsClosed
	}

	timeRecord.UpdatedAt = time.Now()
	// TODO: implement extended validators

	err = timeRecordService.repo.Update(ctx, timeRecord)

	return timeRecord, err
}

func (timeRecordService *TimeRecordService) CloseByTaskID(ctx context.Context, taskID uint64) error {
	searchResult, err := timeRecordService.getActiveTimeRecordsByTaskId(ctx, taskID)
	if err != nil {
		return err
	}
	if len(*searchResult) > 1 {
		return ErrTimeRecordMultipleActive
	}
	for _, timeRecord := range *searchResult {
		now := time.Now()
		timeRecord.EndTime = &now
		timeRecord.IsClosed = true
		timeRecord.UpdatedAt = now
		err = timeRecordService.repo.Update(ctx, &timeRecord)
	}
	return nil
}

func (timeRecordService *TimeRecordService) Delete(ctx context.Context, id uint64) error {
	timeRecord, err := timeRecordService.GetByID(ctx, id)
	if err != nil {
		return err
	}
	return timeRecordService.repo.Delete(ctx, timeRecord)
}

func (timeRecordService *TimeRecordService) createTimeRecordValidate(
	ctx context.Context,
	timeRecord *model.TimeRecord,
) error {
	searchResult, err := timeRecordService.getActiveTimeRecordsByTaskId(ctx, timeRecord.TaskID)
	if err != nil {
		return err
	}
	if len(*searchResult) > 0 {
		return ErrTimeRecordCreationFailed
	}
	return nil
}

func (timeRecordService *TimeRecordService) getActiveTimeRecordsByTaskId(
	ctx context.Context,
	taskID uint64,
) (*[]model.TimeRecord, error) {
	filters := []gormquery.FilterGroup{
		gormquery.NewFilterGroup(
			gormquery.NewFilter("task_id", "=", taskID),
			gormquery.NewFilter("is_closed", "=", false),
		),
	}
	return timeRecordService.repo.GetFilteredTimeRecords(ctx, filters, nil)
}
