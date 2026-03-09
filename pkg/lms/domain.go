package lms

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type LMS struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	Name      string    `gorm:"unique;not null"`
	Path      string    `gorm:"not null"`
	Status    string    `gorm:"default:'active'"`
	Notes     string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (lms *LMS) BeforeCreate(tx *gorm.DB) (err error) {
	if lms.ID == uuid.Nil {
		lms.ID = uuid.New()
	}
	return
}

type VersionUpdateEvent struct {
	LMSID   string `json:"lms_id"`
	Name    string `json:"name"`
	Version string `json:"version"`
}
