package model

import (
	"time"

	"github.com/google/uuid"
)

type Project struct {
	// FIXME id to int
	ID        uint64    `gorm:"type:uuid;primaryKey" json:"id"`
	Name      string    `gorm:"not null" json:"name"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
