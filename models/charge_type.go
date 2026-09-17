package models

import "time"

// ChargeType is one kind of recurring additional charge an institution
// offers besides the grade-based monthly tuition (mensualidad) — e.g.
// "Transporte", "Alimentación". Defined once per institution; a student
// opts into a subset of these via StudentBillingCharge, which may
// override DefaultAmount per student.
type ChargeType struct {
	ID            uint   `gorm:"primaryKey" json:"id"`
	InstitutionID uint   `gorm:"column:institution_id;not null;index;uniqueIndex:idx_charge_types_institution_name" json:"institution_id"`
	Name          string `gorm:"column:name;size:100;not null;uniqueIndex:idx_charge_types_institution_name" json:"name"`
	// DefaultAmount is whole Colombian pesos — the starting price shown
	// when a student opts into this charge, editable per student from there.
	DefaultAmount int64 `gorm:"column:default_amount;not null" json:"default_amount"`
	IsActive      bool  `gorm:"column:is_active;not null;default:true" json:"is_active"`

	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (ChargeType) TableName() string {
	return "charge_types"
}
