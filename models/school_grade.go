package models

import "time"

// SchoolGrade is one grade/level a school offers (e.g. "Noveno"), with
// the tuition amount charged when a Student assigned to it is enrolled
// for a school year — see EnrollmentService.Register, which looks this
// value up instead of generating a placeholder amount. Institution-scoped
// like Role/Requirement: each school manages its own catalog and pricing.
type SchoolGrade struct {
	ID            uint         `gorm:"primaryKey" json:"id"`
	InstitutionID uint         `gorm:"column:institution_id;not null;index;uniqueIndex:idx_school_grades_institution_name" json:"institution_id"`
	Institution   *Institution `gorm:"foreignKey:InstitutionID" json:"-"`
	Name          string       `gorm:"column:name;size:100;not null;uniqueIndex:idx_school_grades_institution_name" json:"name"`
	// TuitionValue is whole Colombian pesos (no subunit in practice),
	// matching Enrollment.Amount's own convention.
	TuitionValue int64 `gorm:"column:tuition_value;not null" json:"tuition_value"`
	// IsActive gates whether the grade can be newly assigned to a student
	// (see SchoolGradeRepository.ListAll, used by the "assign grade"
	// picker) — disabling one never touches students already assigned to
	// it, nor their past enrollments.
	IsActive bool `gorm:"column:is_active;not null;default:true" json:"is_active"`

	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (SchoolGrade) TableName() string {
	return "school_grades"
}
