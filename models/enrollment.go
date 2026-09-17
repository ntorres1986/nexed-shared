package models

import "time"

type EnrollmentStatus string

const (
	EnrollmentStatusPending EnrollmentStatus = "PENDING"
	EnrollmentStatusPaid    EnrollmentStatus = "PAID"
)

// Enrollment represents one Student's matrícula for a given school year —
// created once per (StudentID, SchoolYear), see the unique index. Amount
// is the institution's fixed Institution.EnrollmentFeeAmount at the time
// of registration (matrícula does not vary by grade); DueDate is likewise
// copied from Institution.EnrollmentFeeDueDate at that moment, so it stays
// accurate even if the institution's configured due date changes in a
// later year. Paid via PlacetoPay — see PaymentService/PaymentTransaction.
type Enrollment struct {
	ID            uint     `gorm:"primaryKey" json:"id"`
	InstitutionID uint     `gorm:"column:institution_id;not null;index" json:"institution_id"`
	StudentID     uint     `gorm:"column:student_id;not null;uniqueIndex:idx_enrollment_student_year" json:"student_id"`
	Student       *Student `gorm:"foreignKey:StudentID" json:"student,omitempty"`

	SchoolYear  int    `gorm:"column:school_year;not null;uniqueIndex:idx_enrollment_student_year" json:"school_year"`
	Description string `gorm:"column:description;size:255;not null" json:"description"`
	// Whole Colombian pesos (no subunit in practice).
	Amount  int64            `gorm:"column:amount;not null" json:"amount"`
	DueDate *time.Time       `gorm:"column:due_date" json:"due_date,omitempty"`
	Status  EnrollmentStatus `gorm:"column:status;size:20;not null" json:"status"`
	// ReportedAt is set once a reviewer (ADMINISTRATIVO or a role granted
	// enrollments.approve_requirements) manually confirms this student as
	// formally matriculado from the grade requirement-review grid — nil
	// means not yet confirmed. Independent of Status, which is only about
	// the matrícula fee payment: a fully-paid enrollment isn't
	// automatically "reported", and a reviewer may report one before every
	// document is approved. Setting it also emails the parent — see
	// RequirementReviewService.MarkReported.
	ReportedAt *time.Time `gorm:"column:reported_at" json:"reported_at,omitempty"`

	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (Enrollment) TableName() string {
	return "enrollments"
}
