package models

import "time"

// StudentBillingPlan is one student's monthly-billing configuration for a
// school year — how many installments (typically 10), the monthly
// tuition amount (defaults from SchoolGrade.MonthlyFeeAmount at creation,
// editable per student), and which additional ChargeTypes apply (see
// StudentBillingCharge). Configured from the student form when the admin
// assigns a grade — see StudentBillingPlanService. Saving/changing a plan
// (regenerates) the actual MonthlyDue rows, but only ones still PENDING —
// already-paid history is never altered.
type StudentBillingPlan struct {
	ID            uint     `gorm:"primaryKey" json:"id"`
	InstitutionID uint     `gorm:"column:institution_id;not null;index" json:"institution_id"`
	StudentID     uint     `gorm:"column:student_id;not null;uniqueIndex:idx_billing_plan_student_year" json:"student_id"`
	Student       *Student `gorm:"foreignKey:StudentID" json:"-"`
	SchoolYear    int      `gorm:"column:school_year;not null;uniqueIndex:idx_billing_plan_student_year" json:"school_year"`

	InstallmentCount int   `gorm:"column:installment_count;not null" json:"installment_count"`
	MonthlyAmount    int64 `gorm:"column:monthly_amount;not null" json:"monthly_amount"`

	Charges []StudentBillingCharge `gorm:"foreignKey:BillingPlanID" json:"charges,omitempty"`

	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (StudentBillingPlan) TableName() string {
	return "student_billing_plans"
}

// StudentBillingCharge is one ChargeType a student's billing plan opted
// into, with its own (possibly overridden) monthly Amount and its own
// installment schedule — independent from the plan's base
// StudentBillingPlan.InstallmentCount, since an additional charge (e.g.
// "Transporte") may be a one-time fee (InstallmentCount 1) or span a
// different number of months than the mensualidad, starting in a
// different month too.
type StudentBillingCharge struct {
	ID            uint        `gorm:"primaryKey" json:"id"`
	BillingPlanID uint        `gorm:"column:billing_plan_id;not null;index;uniqueIndex:idx_billing_charge_plan_type" json:"billing_plan_id"`
	ChargeTypeID  uint        `gorm:"column:charge_type_id;not null;uniqueIndex:idx_billing_charge_plan_type" json:"charge_type_id"`
	ChargeType    *ChargeType `gorm:"foreignKey:ChargeTypeID" json:"charge_type,omitempty"`
	// Amount is whole Colombian pesos — defaults from
	// ChargeType.DefaultAmount when added, editable per student.
	Amount int64 `gorm:"column:amount;not null" json:"amount"`
	// InstallmentCount is this charge's own number of cuotas — 1 for a
	// single one-off charge, or however many months it recurs for.
	InstallmentCount int `gorm:"column:installment_count;not null;default:1" json:"installment_count"`
	// StartMonth (1-12) is the calendar month the charge's first
	// installment falls due in; subsequent installments follow one per
	// month from there (same day-of-month as Institution.MonthlyDueDay).
	StartMonth int `gorm:"column:start_month;not null;default:2" json:"start_month"`
}

func (StudentBillingCharge) TableName() string {
	return "student_billing_charges"
}
