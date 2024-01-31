package types

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Base is the base struct of "tables"
// Struct to common fields
type Base struct {
	ID         uuid.UUID      `json:"id" gorm:"primary_key"`
	Created_At time.Time      `json:"created_at" gorm:"autoCreateTime"`
	Updated_At time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	Deleted_At gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}

// BeforeCreate will set a UUID rather than numeric ID.
func (base *Base) BeforeCreate(tx *gorm.DB) error {
	uuid := uuid.New()

	base.ID = uuid

	return nil
}
