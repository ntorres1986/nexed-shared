package models

import "time"

// TemplateVariable is the fixed, safe catalog of values that can be
// injected into a requirement template — deliberately NOT free text, so a
// template can never reference arbitrary/unintended data.
type TemplateVariable string

const (
	TemplateVarStudentFullName   TemplateVariable = "STUDENT_FULL_NAME"
	TemplateVarStudentDocumentID TemplateVariable = "STUDENT_DOCUMENT_ID"
	TemplateVarStudentGrade      TemplateVariable = "STUDENT_GRADE"
	TemplateVarParentFullName    TemplateVariable = "PARENT_FULL_NAME"
	TemplateVarParentEmail       TemplateVariable = "PARENT_EMAIL"
	TemplateVarInstitutionName   TemplateVariable = "INSTITUTION_NAME"
	TemplateVarSchoolYear        TemplateVariable = "SCHOOL_YEAR"
	TemplateVarCurrentDate       TemplateVariable = "CURRENT_DATE"
)

// RequirementTemplate is the uploaded PDF (a school's own document, e.g.
// "Contrato de matrícula", with blank lines meant to be filled by hand)
// attached to one Requirement. Only a Requirement of type
// RequirementTypeForm ever has one — an UPLOAD-type requirement is
// something the PARENT uploads, not the school. One-to-one with
// Requirement.
type RequirementTemplate struct {
	ID               uint                       `gorm:"primaryKey" json:"id"`
	InstitutionID    uint                       `gorm:"column:institution_id;not null;index" json:"institution_id"`
	RequirementID    uint                       `gorm:"column:requirement_id;not null;uniqueIndex" json:"requirement_id"`
	FilePath         string                     `gorm:"column:file_path;size:255;not null" json:"file_path"`
	OriginalFilename string                     `gorm:"column:original_filename;size:255;not null" json:"original_filename"`
	PageCount        int                        `gorm:"column:page_count;not null" json:"page_count"`
	Fields           []RequirementTemplateField `gorm:"foreignKey:TemplateID" json:"fields,omitempty"`
	CreatedAt        time.Time                  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt        time.Time                  `gorm:"column:updated_at" json:"updated_at"`
}

func (RequirementTemplate) TableName() string {
	return "requirement_templates"
}

// RequirementTemplateField is one variable placed at a specific position
// on a specific page of a RequirementTemplate's PDF — the school admin
// places these visually (see requirement_template_handler.go), clicking
// where each blank line falls on the original document. X/Y are stored as
// fractions of the page's width/height (0..1, origin top-left, y growing
// downward — matching how a browser click on a rendered PDF page is
// captured) rather than absolute points, so placement stays correct
// regardless of the zoom level used in the editor; the point conversion
// happens at generation time against that page's actual dimensions (see
// DocumentGenerationService).
type RequirementTemplateField struct {
	ID         uint             `gorm:"primaryKey" json:"id"`
	TemplateID uint             `gorm:"column:template_id;not null;index" json:"template_id"`
	Variable   TemplateVariable `gorm:"column:variable;size:40;not null" json:"variable"`
	PageNumber int              `gorm:"column:page_number;not null" json:"page_number"`
	XRatio     float64          `gorm:"column:x_ratio;not null" json:"x_ratio"`
	YRatio     float64          `gorm:"column:y_ratio;not null" json:"y_ratio"`
	FontSize   int              `gorm:"column:font_size;not null" json:"font_size"`
	CreatedAt  time.Time        `gorm:"column:created_at" json:"created_at"`
	UpdatedAt  time.Time        `gorm:"column:updated_at" json:"updated_at"`
}

func (RequirementTemplateField) TableName() string {
	return "requirement_template_fields"
}
