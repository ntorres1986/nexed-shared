package models

import "time"

type MonthlyDueStatus string

const (
	MonthlyDueStatusPending MonthlyDueStatus = "PENDING"
	MonthlyDueStatusPaid    MonthlyDueStatus = "PAID"
)

// MonthlyDue is one real, individually-payable installment generated from
// a StudentBillingPlan — either the base monthly tuition (ChargeTypeID
// nil) or one installment of a specific additional ChargeType. Each
// "stream" (the base tuition, or each opted-in charge) has its own
// InstallmentNumber sequence 1..InstallmentCount, so a student with the
// base tuition plus one extra charge over 10 installments has 20 rows,
// not 10. Paid via PlacetoPay exactly like Enrollment — see
// PaymentTransaction.MonthlyDueID and PaymentService.
type MonthlyDue struct {
	ID            uint     `gorm:"primaryKey" json:"id"`
	InstitutionID uint     `gorm:"column:institution_id;not null;index" json:"institution_id"`
	StudentID     uint     `gorm:"column:student_id;not null;index;uniqueIndex:idx_monthly_due_stream,priority:1" json:"student_id"`
	Student       *Student `gorm:"foreignKey:StudentID" json:"student,omitempty"`
	SchoolYear    int      `gorm:"column:school_year;not null;uniqueIndex:idx_monthly_due_stream,priority:2" json:"school_year"`

	// InstallmentNumber is 1-based within its own stream (see ChargeTypeID).
	InstallmentNumber int `gorm:"column:installment_number;not null;uniqueIndex:idx_monthly_due_stream,priority:4" json:"installment_number"`
	// ChargeTypeID is nil for the base monthly tuition stream, or the
	// specific additional charge this installment belongs to.
	ChargeTypeID *uint       `gorm:"column:charge_type_id;uniqueIndex:idx_monthly_due_stream,priority:3" json:"charge_type_id,omitempty"`
	ChargeType   *ChargeType `gorm:"foreignKey:ChargeTypeID" json:"charge_type,omitempty"`

	DueDate     time.Time        `gorm:"column:due_date;not null" json:"due_date"`
	Description string           `gorm:"column:description;size:255;not null" json:"description"`
	Amount      int64            `gorm:"column:amount;not null" json:"amount"`
	Status      MonthlyDueStatus `gorm:"column:status;size:20;not null;index" json:"status"`
	// DueDateOverridden is true once an admin manually edited this
	// installment's DueDate (see StudentBillingPlanService.UpdateDueDate) —
	// a later StudentBillingPlan save then leaves DueDate alone for this
	// row instead of recomputing it, so the manual override survives
	// unrelated plan edits (price, other installments, other charges).
	DueDateOverridden bool `gorm:"column:due_date_overridden;not null;default:false" json:"due_date_overridden"`

	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (MonthlyDue) TableName() string {
	return "monthly_dues"
}
