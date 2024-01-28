package types

import (
	"time"

	"github.com/google/uuid"
)

// Base is the base struct of "tables"
// Struct to common fields
type Base struct {
	ID         uuid.UUID  `json:"id" bson:"_id,omitempty"`
	Active     bool       `json:"-" bson:"active"`
	Created_At time.Time  `json:"created_at" bson:"created_at"`
	Updated_At time.Time  `json:"updates_at" bson:"updates_at"`
	Deletes_At *time.Time `json:"deleted_at,omitempty" bson:"deleted_at,omitempty"`
}
