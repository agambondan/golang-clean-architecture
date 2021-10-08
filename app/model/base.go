package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type BaseInt struct {
	ID *int64 `json:"id,omitempty"`
	BaseDate
}

type BaseUUID struct {
	ID *uuid.UUID `json:"id,omitempty" gorm:"primaryKey;unique;type:varchar(36);not null" format:"uuid"` // model ID
	BaseDate
}

type BaseDate struct {
	CreatedAt *time.Time     `json:"created_at,omitempty" gorm:"type:timestamptz" format:"date-time" swaggertype:"string"` // created at automatically inserted on post
	UpdatedAt *time.Time     `json:"updated_at,omitempty" gorm:"type:timestamptz" format:"date-time" swaggertype:"string"` // updated at automatically changed on put or add on post
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index" swaggerignore:"true"`
}

func (b *BaseUUID) BeforeCreate(db *gorm.DB) error {
	if b.ID == nil {
		newUUID, _ := uuid.NewUUID()
		b.ID = &newUUID
	}
	return nil
}

func (b *BaseDate) BeforeCreate(db *gorm.DB) error {
	now := time.Now()
	b.CreatedAt = &now
	b.UpdatedAt = &now
	return nil
}

// BeforeUpdate Data
func (b *BaseDate) BeforeUpdate(tx *gorm.DB) error {
	now := time.Now()
	b.UpdatedAt = &now
	return nil
}
