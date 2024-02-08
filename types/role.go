package types

import (
	"encoding/json"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type Role struct {
	Base
	Name        string      `json:"name" gorm:"type:varchar(25)"`
	Description string      `json:"description"`
	Tag         string      `json:"tag" gorm:"unique"`
	App         Application `json:"-" gorm:"foreignKey:AppID"`
	AppID       uuid.UUID   `json:"application_id" gorm:"index"`
	Permission  Permission  `json:"permission" gorm:"type:text"`
}

type Permission struct {
	Config []PermissionConfig `json:"config"`
}

type PermissionConfig struct {
	Path           string          `json:"path"`
	AllowedMethods map[string]bool `json:"allowed_methods"`
}

// Scan Unmarshal custom JSON type
func (p *Permission) Scan(value interface{}) error {
	data, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(data, &p)
}

func (Permission) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	// use field.Tag, field.TagSettings gets field's tags
	// checkout https://github.com/go-gorm/gorm/blob/master/schema/field.go for all options

	// returns different database type based on driver name
	switch db.Dialector.Name() {
	case "mysql", "sqlite":
		return "json"
	case "postgres":
		return "jsonb"
	}
	return ""
}
