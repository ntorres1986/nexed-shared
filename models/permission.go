package models

import "time"

// Permission is a granular capability, named "<module>.<action>"
// (e.g. "payments.view"). New permissions can be added purely as data
// (via seed or a future admin UI) without touching structural code.
type Permission struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"column:name;size:100;not null;uniqueIndex" json:"name"`
	Description string    `gorm:"column:description;size:255" json:"description"`
	Module      string    `gorm:"column:module;size:50;not null;index" json:"module"`
	Action      string    `gorm:"column:action;size:50;not null" json:"action"`
	CreatedAt   time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (Permission) TableName() string {
	return "permissions"
}
