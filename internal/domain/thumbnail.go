package domain

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Thumbnail struct {
	ID      uuid.UUID `gorm:"type:uuid;primaryKey"`
	ImageID uuid.UUID `gorm:"type:uuid;not null;index:uniq_file_size,unique"`
	Size    string    `gorm:"not null;index:uniq_file_size,unique"`
	Path    string    `gorm:"not null"`
	Type    string    `gorm:"not null"`
}

func (f *Thumbnail) BeforeCreate(tx *gorm.DB) (err error) {
	if f.ID == uuid.Nil {
		f.ID = uuid.New()
	}
	return
}
