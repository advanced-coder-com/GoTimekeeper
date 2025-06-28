package repository

import (
	"context"
	"fmt"
	"github.com/advanced-coder-com/go-timekeeper/internal/db"
	"github.com/advanced-coder-com/go-timekeeper/internal/gormquery"
	"github.com/advanced-coder-com/go-timekeeper/internal/model"
	"gitlab.com/tozd/go/errors"
	"gorm.io/gorm"
)

type TimeRecordRepository interface {
	Create(ctx context.Context, timeRecord *model.TimeRecord) error
	GetByID(ctx context.Context, id uint64) (*model.TimeRecord, error)
	GetByTaskID(ctx context.Context, taskID uint64) (*[]model.TimeRecord, error)
	GetFilteredTimeRecords(
		ctx context.Context,
		filters []gormquery.FilterGroup,
		options *gormquery.QueryOptions,
	) (*[]model.TimeRecord, error)
	Update(ctx context.Context, timeRecord *model.TimeRecord) error
	Delete(ctx context.Context, timeRecord *model.TimeRecord) error
}

type timeRecordRepository struct {
	database *gorm.DB
}

const timeRecordRepoErrorPrefix = "TimeRecordRepository"

func NewTimeRecordRepository() TimeRecordRepository {
	return &timeRecordRepository{database: db.Get()}
}

func (timeRecordRepo *timeRecordRepository) Create(ctx context.Context, timeRecord *model.TimeRecord) error {
	err := timeRecordRepo.database.WithContext(ctx).Create(timeRecord).Error
	if err != nil {
		err = fmt.Errorf("%s create time record failed: %w", timeRecordRepoErrorPrefix, err)
	}
	return err
}

func (timeRecordRepo *timeRecordRepository) GetByID(ctx context.Context, id uint64) (*model.TimeRecord, error) {
	var timeRecord model.TimeRecord
	filters := []gormquery.FilterGroup{
		gormquery.NewFilterGroup(
			gormquery.NewFilter("id", "=", id),
		),
	}
	query := timeRecordRepo.database.WithContext(ctx).Model(&model.TimeRecord{})
	query = gormquery.ApplyFilters(query, filters)
	if err := query.First(&timeRecord).Error; err != nil {
		return nil, fmt.Errorf("%s find time record by id failed: %w", timeRecordRepoErrorPrefix, err)
	}
	return &timeRecord, nil
}

func (timeRecordRepo *timeRecordRepository) GetByTaskID(ctx context.Context, taskID uint64) (*[]model.TimeRecord, error) {
	var timeRecords []model.TimeRecord
	filters := []gormquery.FilterGroup{
		gormquery.NewFilterGroup(
			gormquery.NewFilter("task_id", "=", taskID),
		),
	}
	query := timeRecordRepo.database.WithContext(ctx).Model(&model.TimeRecord{})
	query = gormquery.ApplyFilters(query, filters)
	if err := query.Find(&timeRecords).Error; err != nil {
		err = fmt.Errorf("%s find time records by task id failed: %w", timeRecordRepoErrorPrefix, err)
		return nil, err
	}

	return &timeRecords, nil
}

func (timeRecordRepo *timeRecordRepository) GetFilteredTimeRecords(
	ctx context.Context,
	filters []gormquery.FilterGroup,
	options *gormquery.QueryOptions,
) (*[]model.TimeRecord, error) {
	var timeRecords []model.TimeRecord

	query := timeRecordRepo.database.WithContext(ctx).Model(&model.TimeRecord{})
	query = gormquery.ApplyFilters(query, filters)
	if options != nil {
		query = gormquery.ApplyQueryOptions(query, *options)
	}
	err := query.Find(&timeRecords).Error
	if err != nil {
		err = fmt.Errorf("%s find filtered time records failed: %w", timeRecordRepoErrorPrefix, err)
	}
	return &timeRecords, err
}

func (timeRecordRepo *timeRecordRepository) Update(ctx context.Context, timeRecord *model.TimeRecord) error {
	err := timeRecordRepo.database.WithContext(ctx).Save(timeRecord).Error
	if err != nil {
		err = fmt.Errorf("%s update time record failed: %w", timeRecordRepoErrorPrefix, err)
	}
	return err
}

func (timeRecordRepo *timeRecordRepository) Delete(ctx context.Context, timeRecord *model.TimeRecord) error {
	result := timeRecordRepo.database.
		WithContext(ctx).
		Where("id = ?", timeRecord.ID).
		Delete(&model.TimeRecord{})
	if result.Error != nil {
		return fmt.Errorf("%s delete time record failed: %w", timeRecordRepoErrorPrefix, result.Error)
	}
	if result.RowsAffected == 0 {
		return errors.New(
			fmt.Sprintf(
				"%s delete time record failed: time record you try to delete does not exist",
				timeRecordRepoErrorPrefix,
			),
		)
	}
	return nil
}
