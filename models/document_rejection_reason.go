package models

import "time"

// DocumentRejectionReason is one entry in an institution's configurable
// catalog of reasons a reviewer can give when rejecting a parent-uploaded
// document (see RequirementSubmission.RejectionReasonID in
// nexed-customer-backend) — picked from a grid instead of free text, with
// a "custom reason" escape hatch handled separately by the caller.
type DocumentRejectionReason struct {
	ID            uint   `gorm:"primaryKey" json:"id"`
	InstitutionID uint   `gorm:"column:institution_id;not null;index;uniqueIndex:idx_doc_rejection_reasons_institution_title" json:"institution_id"`
	Title         string `gorm:"column:title;size:150;not null;uniqueIndex:idx_doc_rejection_reasons_institution_title" json:"title"`
	Description   string `gorm:"column:description;size:500" json:"description"`
	IsActive      bool   `gorm:"column:is_active;not null;default:true" json:"is_active"`

	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (DocumentRejectionReason) TableName() string {
	return "document_rejection_reasons"
}
