package models

import "time"

// RequirementType distinguishes a requirement filled out directly on the
// platform ("Documento en línea" — a form/e-signature flow, not built yet)
// from one the parent must upload as a file ("Cargar documento").
type RequirementType string

const (
	RequirementTypeForm   RequirementType = "FORM"
	RequirementTypeUpload RequirementType = "UPLOAD"
)

// RequirementCategory groups a set of Requirements under one heading
// (e.g. "Documentos que debe cargar a la plataforma"). Institution-scoped
// like Role — each institution owns its own independent catalog, cloned
// with sensible defaults when the institution is created (see
// customer-backend/internal/services/requirement_cloning.go).
type RequirementCategory struct {
	ID            uint          `gorm:"primaryKey" json:"id"`
	InstitutionID uint          `gorm:"column:institution_id;not null;index" json:"institution_id"`
	Name          string        `gorm:"column:name;size:150;not null" json:"name"`
	SortOrder     int           `gorm:"column:sort_order;not null" json:"sort_order"`
	Requirements  []Requirement `gorm:"foreignKey:CategoryID" json:"requirements,omitempty"`
	CreatedAt     time.Time     `gorm:"column:created_at" json:"created_at"`
	UpdatedAt     time.Time     `gorm:"column:updated_at" json:"updated_at"`
}

func (RequirementCategory) TableName() string {
	return "requirement_categories"
}

// Requirement is one document/step a Student's Enrollment must satisfy —
// the catalog definition (what's required), not an individual
// enrollment's progress on it (see EnrollmentRequirement).
type Requirement struct {
	ID            uint            `gorm:"primaryKey" json:"id"`
	InstitutionID uint            `gorm:"column:institution_id;not null;index" json:"institution_id"`
	CategoryID    uint            `gorm:"column:category_id;not null;index" json:"category_id"`
	Name          string          `gorm:"column:name;size:150;not null" json:"name"`
	Description   string          `gorm:"column:description;size:150;not null" json:"description"`
	Type          RequirementType `gorm:"column:type;size:20;not null" json:"type"`
	SortOrder     int             `gorm:"column:sort_order;not null" json:"sort_order"`
	CreatedAt     time.Time       `gorm:"column:created_at" json:"created_at"`
	UpdatedAt     time.Time       `gorm:"column:updated_at" json:"updated_at"`
}

func (Requirement) TableName() string {
	return "requirements"
}

// RequirementStatus is one Enrollment's progress on one Requirement.
type RequirementStatus string

const (
	RequirementStatusMissing  RequirementStatus = "FALTA"
	RequirementStatusInReview RequirementStatus = "EN_REVISION"
	RequirementStatusApproved RequirementStatus = "APROBADO"
	RequirementStatusRejected RequirementStatus = "RECHAZADO"
)

// EnrollmentRequirement tracks one Enrollment's progress against one
// Requirement. Rows are created lazily, only once there's something to
// record — an enrollment with no row for a given requirement is
// implicitly RequirementStatusMissing (see
// RequirementService.GetForEnrollment), so nothing needs to eagerly
// initialize one row per requirement when an enrollment is created.
type EnrollmentRequirement struct {
	ID            uint              `gorm:"primaryKey" json:"id"`
	EnrollmentID  uint              `gorm:"column:enrollment_id;not null;uniqueIndex:idx_enrollment_requirement" json:"enrollment_id"`
	RequirementID uint              `gorm:"column:requirement_id;not null;uniqueIndex:idx_enrollment_requirement" json:"requirement_id"`
	Requirement   *Requirement      `gorm:"foreignKey:RequirementID" json:"requirement,omitempty"`
	Status        RequirementStatus `gorm:"column:status;size:20;not null" json:"status"`
	CreatedAt     time.Time         `gorm:"column:created_at" json:"created_at"`
	UpdatedAt     time.Time         `gorm:"column:updated_at" json:"updated_at"`
}

func (EnrollmentRequirement) TableName() string {
	return "enrollment_requirements"
}
