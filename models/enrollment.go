package models

import "time"

type EnrollmentStatus string

const (
	EnrollmentStatusPending EnrollmentStatus = "PENDING"
	EnrollmentStatusPaid    EnrollmentStatus = "PAID"
)

// Enrollment represents one Student's matrícula for a given school year —
// created once per (StudentID, SchoolYear), see the unique index, with a
// placeholder Amount rolled at creation time (no real pricing engine yet)
// and persisted so revisiting the "Realizar matrícula" screen always
// shows the same value. "HACER PAGO" moves Status from PENDING to PAID;
// there is no real payment gateway behind this yet.
type Enrollment struct {
	ID            uint     `gorm:"primaryKey" json:"id"`
	InstitutionID uint     `gorm:"column:institution_id;not null;index" json:"institution_id"`
	StudentID     uint     `gorm:"column:student_id;not null;uniqueIndex:idx_enrollment_student_year" json:"student_id"`
	Student       *Student `gorm:"foreignKey:StudentID" json:"student,omitempty"`

	SchoolYear  int    `gorm:"column:school_year;not null;uniqueIndex:idx_enrollment_student_year" json:"school_year"`
	Description string `gorm:"column:description;size:255;not null" json:"description"`
	// Whole Colombian pesos (no subunit in practice) — see EnrollmentService.randomAmount.
	Amount int64            `gorm:"column:amount;not null" json:"amount"`
	Status EnrollmentStatus `gorm:"column:status;size:20;not null" json:"status"`

	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (Enrollment) TableName() string {
	return "enrollments"
}
