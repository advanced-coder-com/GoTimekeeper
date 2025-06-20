package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type TaskStatus string

const (
	StatusOpened    TaskStatus = "Opened"
	StatusWorkingOn TaskStatus = "Working on"
	StatusClosed    TaskStatus = "Closed"
	DefaultStatus              = StatusOpened
)

type Task struct {
	ID        uint64         `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    uuid.UUID      `gorm:"type:uuid;not null" json:"user_id"`
	ProjectID uint64         `gorm:"type:uuid" json:"project_id,omitempty"`
	Name      string         `gorm:"not null" json:"name"`
	Tags      pq.StringArray `gorm:"type:text[]" json:"tags,omitempty"`
	Status    TaskStatus     `gorm:"type:varchar(20);not null" json:"status"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

func IsValidTaskStatus(inputStatus string) bool {
	switch TaskStatus(inputStatus) {
	case StatusOpened, StatusWorkingOn, StatusClosed:
		return true
	default:
		return false
	}
}
