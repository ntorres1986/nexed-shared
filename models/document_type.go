package models

import "time"

// DocumentType is one entry in an institution's configurable catalog of
// identity document types (cédula de ciudadanía, tarjeta de identidad,
// registro civil, pasaporte, etc.) — shared between Student and User
// (a PADRE's own document), since it's the same real-world concept for
// both.
type DocumentType struct {
	ID            uint   `gorm:"primaryKey" json:"id"`
	InstitutionID uint   `gorm:"column:institution_id;not null;index;uniqueIndex:idx_document_types_institution_name" json:"institution_id"`
	Name          string `gorm:"column:name;size:100;not null;uniqueIndex:idx_document_types_institution_name" json:"name"`
	IsActive      bool   `gorm:"column:is_active;not null;default:true" json:"is_active"`

	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (DocumentType) TableName() string {
	return "document_types"
}
