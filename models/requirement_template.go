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

// RequirementTemplate is one version of a FORM requirement's document —
// rich text (authored with an inline WYSIWYG editor, variables inserted as
// chips) scoped to either every grade or a specific subset of them. A
// requirement can have several templates at once (e.g. one "general" one
// plus one for a particular grade whose conditions differ) — see
// RequirementTemplateRepository.FindBestMatch for how the right one is
// picked for a given student at generation time. Unlike the old
// coordinate-stamped-PDF design, the document itself is generated fresh
// from BodyHTML each time (see template_pdf.go), repeating the
// institution's own Institution.LetterheadPath as a header on every page.
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
	// BodyHTML is sanitized server-side before it's ever saved (see
	// RequirementTemplateService's bluemonday policy) — it's rendered
	// straight into a legal document a parent signs, not just displayed.
	BodyHTML string `gorm:"column:body_html;type:text;not null" json:"body_html"`
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
