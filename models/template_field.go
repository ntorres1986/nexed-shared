package models

import "time"

// TemplateField places one variable's value at a fixed position over a
// PDF-sourced RequirementTemplate (see RequirementTemplate.SourceType).
// Only meaningful when the owning template's SourceType is PDF.
type TemplateField struct {
	ID         uint `gorm:"primaryKey" json:"id"`
	TemplateID uint `gorm:"column:template_id;not null;index" json:"template_id"`
	// Variable is one of the TemplateVariable constants — the same catalog
	// used by the TEXT source type's chips.
	Variable TemplateVariable `gorm:"column:variable;size:50;not null" json:"variable"`
	// Page is 1-indexed, matching the uploaded PDF's own page numbering.
	Page int `gorm:"column:page;not null" json:"page"`
	// XFraction/YFraction are 0..1, relative to that page's own width/height,
	// measured from the top-left corner — the same convention an HTML
	// canvas uses, so the frontend editor's coordinates apply unchanged.
	XFraction float64 `gorm:"column:x_fraction;not null" json:"x_fraction"`
	YFraction float64 `gorm:"column:y_fraction;not null" json:"y_fraction"`
	// FontSize is in points; the service defaults it (e.g. 11) if unset.
	FontSize  float64   `gorm:"column:font_size;not null" json:"font_size"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (TemplateField) TableName() string {
	return "template_fields"
}
