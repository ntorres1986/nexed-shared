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
	RequirementsNoticeHTML string    `gorm:"column:requirements_notice_html;type:text" json:"requirements_notice_html"`
	CreatedAt              time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt              time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (Institution) TableName() string {
	return "institutions"
}
