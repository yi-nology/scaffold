package model

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

// StringArray custom type for storing string arrays as JSON
type StringArray []string

func (s *StringArray) Scan(value interface{}) error {
	if value == nil {
		*s = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		str, ok := value.(string)
		if !ok {
			*s = nil
			return nil
		}
		bytes = []byte(str)
	}
	return json.Unmarshal(bytes, s)
}

func (s StringArray) Value() (driver.Value, error) {
	if s == nil {
		return nil, nil
	}
	return json.Marshal(s)
}

// TemplateModel database model for templates
type TemplateModel struct {
	ID          string      `gorm:"primaryKey;size:100" json:"id"`
	Name        string      `gorm:"size:255;not null;index" json:"name"`
	Description string      `gorm:"type:text" json:"description"`
	Author      string      `gorm:"size:100" json:"author"`
	Version     string      `gorm:"size:50" json:"version"`
	Repository  string      `gorm:"size:500" json:"repository"`
	LocalPath   string      `gorm:"size:500" json:"local_path"`
	Tags        StringArray `gorm:"type:json" json:"tags"`
	CreatedAt   time.Time   `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time   `gorm:"autoUpdateTime" json:"updated_at"`
}

func (TemplateModel) TableName() string {
	return "templates"
}
