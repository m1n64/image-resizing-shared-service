package domain

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type ImageStatus string

const (
	StatusPending    ImageStatus = "pending"
	StatusProcessing ImageStatus = "processing"
	StatusReady      ImageStatus = "ready"
	StatusError      ImageStatus = "error"
)

type Image struct {
	ID           uuid.UUID   `gorm:"type:uuid;primaryKey"`
	Filename     string      `gorm:"not null"`
	MIME         string      `gorm:"not null"`
	Size         int64       `gorm:"not null"`
	OriginalPath string      `gorm:"not null"`
	Status       ImageStatus `gorm:"not null;default:'pending'"`
	ErrorMessage *string     `gorm:""`
	WebPPath     string      `gorm:""`
	CreatedAt    time.Time   `gorm:"autoCreateTime"`
	Thumbnails   []Thumbnail `gorm:"foreignKey:ImageID;constraint:OnDelete:CASCADE"`
}

func (f *Image) BeforeCreate(tx *gorm.DB) (err error) {
	if f.ID == uuid.Nil {
		f.ID = uuid.New()
	}
	return
}
