package models

import "time"

// Institution represents a school/tenant. Every institution-scoped table
// carries an institution_id that must trace back to a row here.
type Institution struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	Name     string `gorm:"column:name;size:150;not null" json:"name"`
	Code     string `gorm:"column:code;size:50;not null;uniqueIndex" json:"code"`
	Logo     string `gorm:"column:logo;size:255" json:"logo"`
	IsActive bool   `gorm:"column:is_active;not null" json:"is_active"`
	// Subdomain is this institution's slice of the multi-tenant deployment
	// (e.g. "colegio1" for colegio1.tudominio.com) — see
	// internal/middleware's Host-based institution resolution, which every
	// request goes through before auth. Nullable (not every institution
	// has one set, e.g. local dev) and unique when set; MySQL's unique
	// index allows any number of NULLs, so this stays optional per-row.
	Subdomain *string `gorm:"column:subdomain;size:63;uniqueIndex" json:"subdomain,omitempty"`
	// RequirementsNoticeHTML is the sanitized rich-text warning shown at the
	// top of a student's requirements checklist (e.g. "your enrollment
	// isn't final until every document is approved"). Institution-configurable
	// from Opciones → General; empty means the frontend falls back to its
	// own default copy.
	RequirementsNoticeHTML string `gorm:"column:requirements_notice_html;type:text" json:"requirements_notice_html"`
	// PrimaryColor/SecondaryColor are the institution's brand colors (hex,
	// e.g. "#7c3aed"), set from nexed-admin and applied across
	// nexed-customer-frontend (theme, sidebar, buttons). Empty means "use
	// the platform default" — that fallback lives client-side, not here.
	PrimaryColor   string `gorm:"column:primary_color;size:7" json:"primary_color"`
	SecondaryColor string `gorm:"column:secondary_color;size:7" json:"secondary_color"`
	// LetterheadPath is the institution's single letterhead image, stamped
	// as a repeating header on every page of a generated requirement
	// template document (see RequirementTemplate) — one per institution,
	// reused across all of its templates. Empty means templates render
	// without a header band.
	LetterheadPath string `gorm:"column:letterhead_path;size:255" json:"letterhead_path"`
	// EnrollmentFeeAmount is the matrícula fee — a single fixed amount for
	// every student regardless of grade (unlike SchoolGrade.MonthlyFeeAmount,
	// which does vary by grade). EnrollmentFeeDueDate is copied onto each
	// Enrollment at registration time rather than read live, so a later
	// change here doesn't retroactively alter what an already-registered
	// family was told — see Enrollment.DueDate.
	EnrollmentFeeAmount  int64      `gorm:"column:enrollment_fee_amount;not null" json:"enrollment_fee_amount"`
	EnrollmentFeeDueDate *time.Time `gorm:"column:enrollment_fee_due_date" json:"enrollment_fee_due_date,omitempty"`
	// MonthlyDueDay is the day of the month (1-28, to stay valid in every
	// month including February) that every MonthlyDue installment falls
	// due on — see StudentBillingPlanService's due-date generation, which
	// reads this instead of a hardcoded day. Defaults to 5.
	MonthlyDueDay int       `gorm:"column:monthly_due_day;not null;default:5" json:"monthly_due_day"`
	CreatedAt     time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt     time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (Institution) TableName() string {
	return "institutions"
}
