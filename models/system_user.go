package models

import "time"

// SystemUser is a platform administrator (what the `admin/` frontend
// authenticates as) — deliberately kept in its own table, fully separate
// from `users` (institution accounts: PADRE/DOCENTE/ADMINISTRATIVO).
//
// Mixing platform staff with tenant data in the same table is a common
// multi-tenant anti-pattern: it couples an account type that must never be
// institution-scoped with one that always must be, and makes it easy for
// a bug in institution-scoping logic to accidentally expose or be
// affected by admin accounts. A SystemUser has no institution_id and no
// role_id — every system user is, by definition, a fully-privileged
// platform administrator (there is currently only one tier of platform
// staff); RBAC (roles/permissions) exists only for institution accounts.
type SystemUser struct {
	ID                 uint      `gorm:"primaryKey" json:"id"`
	FirstName          string    `gorm:"column:first_name;size:100;not null" json:"first_name"`
	LastName           string    `gorm:"column:last_name;size:100;not null" json:"last_name"`
	Email              string    `gorm:"column:email;size:150;not null;uniqueIndex" json:"email"`
	Password           string    `gorm:"column:password;size:255;not null" json:"-"`
	IsActive           bool      `gorm:"column:is_active;not null" json:"is_active"`
	MustChangePassword bool      `gorm:"column:must_change_password;not null" json:"must_change_password"`
	CreatedAt          time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt          time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (SystemUser) TableName() string {
	return "system_users"
}

func (u SystemUser) FullName() string {
	return u.FirstName + " " + u.LastName
}
