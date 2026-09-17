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

// TemplateSourceType selects how a RequirementTemplate's document is
// produced: TEXT flows BodyHTML top-to-bottom on a fresh page (see
// template_pdf.go's RenderTemplatePDF); PDF stamps variable values onto
// fixed x/y positions over an admin-uploaded PDF (RenderTemplatePDFFromSource),
// letting the institution design the layout in an external tool and just
// mark where each value goes.
type TemplateSourceType string

const (
	TemplateSourceTypeText TemplateSourceType = "TEXT"
	TemplateSourceTypePDF  TemplateSourceType = "PDF"
)

// RequirementTemplate is one version of a FORM requirement's document,
// scoped to either every grade or a specific subset of them. A requirement
// can have several templates at once (e.g. one "general" one plus one for a
// particular grade whose conditions differ) — see
// RequirementTemplateRepository.FindBestMatch for how the right one is
// picked for a given student at generation time.
type RequirementTemplate struct {
	ID            uint `gorm:"primaryKey" json:"id"`
	InstitutionID uint `gorm:"column:institution_id;not null;index" json:"institution_id"`
	// RequirementID is intentionally NOT unique — a requirement can now
	// have multiple templates (one per grade-scope combination).
	RequirementID uint `gorm:"column:requirement_id;not null;index" json:"requirement_id"`
	// Name is the admin-facing label distinguishing this template from
	// any siblings on the same requirement, e.g. "Contrato de matrícula -
	// Preescolar".
	Name string `gorm:"column:name;size:150;not null" json:"name"`
	// SourceType picks which of BodyHTML or (PDFPath + Fields) is used to
	// generate this template's document.
	SourceType TemplateSourceType `gorm:"column:source_type;size:10;not null;default:'TEXT'" json:"source_type"`
	// BodyHTML is sanitized server-side before it's ever saved (see
	// RequirementTemplateService's bluemonday policy) — it's rendered
	// straight into a legal document a parent signs, not just displayed.
	// Only meaningful when SourceType is TEXT.
	BodyHTML string `gorm:"column:body_html;type:text;not null" json:"body_html"`
	// PDFPath is the storage path (FileStorage.Save's return value) of the
	// admin-uploaded base PDF. Only set when SourceType is PDF.
	PDFPath string `gorm:"column:pdf_path;size:255" json:"pdf_path,omitempty"`
	// PDFPageCount is captured once at upload time so the frontend editor
	// doesn't need to reparse the PDF just to know how many pages to offer.
	PDFPageCount int `gorm:"column:pdf_page_count" json:"pdf_page_count,omitempty"`
	// Fields are the variable placements over PDFPath. Only populated when
	// SourceType is PDF.
	Fields []TemplateField `gorm:"foreignKey:TemplateID" json:"fields,omitempty"`
	// AppliesToAllGrades, when true, makes this the fallback template used
	// for any student whose grade isn't explicitly covered by one of the
	// requirement's other templates. Grades is meaningless when this is
	// true (and should be empty).
	AppliesToAllGrades bool          `gorm:"column:applies_to_all_grades;not null" json:"applies_to_all_grades"`
	Grades             []SchoolGrade `gorm:"many2many:requirement_template_grades;" json:"grades,omitempty"`
	CreatedAt          time.Time     `gorm:"column:created_at" json:"created_at"`
	UpdatedAt          time.Time     `gorm:"column:updated_at" json:"updated_at"`
}

func (RequirementTemplate) TableName() string {
	return "requirement_templates"
}
