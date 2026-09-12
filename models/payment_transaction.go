package models

import "time"

type PaymentTransactionStatus string

const (
	PaymentStatusPending  PaymentTransactionStatus = "PENDING"
	PaymentStatusApproved PaymentTransactionStatus = "APPROVED"
	PaymentStatusRejected PaymentTransactionStatus = "REJECTED"
	PaymentStatusFailed   PaymentTransactionStatus = "FAILED"
)

// PaymentTransaction is one payment attempt against PlacetoPay/AvalPay
// Checkout for an Enrollment. A parent can retry after a rejected/failed
// attempt (each retry needs its own never-reused Reference — see the WC
// certification guide's field-validation section), so an Enrollment can
// have several of these; the most recent one is the "current" attempt.
// Enrollment.Status flips to PAID the moment any one of them reaches
// APPROVED — see PaymentService.
type PaymentTransaction struct {
	ID           uint        `gorm:"primaryKey" json:"id"`
	EnrollmentID uint        `gorm:"column:enrollment_id;not null;index" json:"enrollment_id"`
	Enrollment   *Enrollment `gorm:"foreignKey:EnrollmentID" json:"-"`

	// Reference is the unique-per-attempt alphanumeric code (max 32 chars)
	// sent to PlacetoPay as payment.reference — never reused across
	// attempts.
	Reference string `gorm:"column:reference;size:32;not null;uniqueIndex" json:"reference"`
	// RequestID is PlacetoPay's own session identifier, returned by
	// CreateSession — used for QuerySession and to match notifications.
	RequestID  int64  `gorm:"column:request_id;index" json:"request_id"`
	ProcessURL string `gorm:"column:process_url;size:500" json:"process_url,omitempty"`

	Status  PaymentTransactionStatus `gorm:"column:status;size:20;not null;index" json:"status"`
	Reason  string                   `gorm:"column:reason;size:50" json:"reason,omitempty"`
	Message string                   `gorm:"column:message;size:255" json:"message,omitempty"`

	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (PaymentTransaction) TableName() string {
	return "payment_transactions"
}
