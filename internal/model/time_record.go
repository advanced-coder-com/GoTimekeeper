package model

import (
	"github.com/google/uuid"
	"time"
)

type TimeRecord struct {
	ID        uint       `gorm:"primaryKey" json:"id"`
	UserID    uuid.UUID  `gorm:"type:uuid;not null" json:"user_id"`
	TaskID    int64      `gorm:"not null;index" json:"task_id"`
	StartTime time.Time  `gorm:"not null" json:"start_time"`
	EndTime   *time.Time `json:"end_time,omitempty"`
	IsClosed  bool       `gorm:"default:false" json:"is_closed"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}
