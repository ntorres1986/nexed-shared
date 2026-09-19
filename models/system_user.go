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
// role_id: RBAC (the roles/permissions catalog) exists only for institution
// accounts. Platform staff instead have a fixed, small Role (see
// SystemRole) — enforced by the admin backend on every route.
type SystemUser struct {
	ID                 uint   `gorm:"primaryKey" json:"id"`
	FirstName          string `gorm:"column:first_name;size:100;not null" json:"first_name"`
	LastName           string `gorm:"column:last_name;size:100;not null" json:"last_name"`
	Email              string `gorm:"column:email;size:150;not null;uniqueIndex" json:"email"`
	Password           string `gorm:"column:password;size:255;not null" json:"-"`
	IsActive           bool   `gorm:"column:is_active;not null" json:"is_active"`
	MustChangePassword bool   `gorm:"column:must_change_password;not null" json:"must_change_password"`
	// Role defaults to SUPER_ADMIN so every account created before roles
	// existed (the original platform administrators) keeps full access.
	Role      SystemRole `gorm:"column:role;size:20;not null;default:SUPER_ADMIN" json:"role"`
	CreatedAt time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time  `gorm:"column:updated_at" json:"updated_at"`
}

// SystemRole is what a platform-staff account may do in the admin panel.
type SystemRole string

const (
	// SystemRoleSuperAdmin can do everything, including managing the other
	// platform users and the permission catalog.
	SystemRoleSuperAdmin SystemRole = "SUPER_ADMIN"
	// SystemRoleAdmin manages institutions and their users.
	SystemRoleAdmin SystemRole = "ADMIN"
	// SystemRoleSupport is read-only.
	SystemRoleSupport SystemRole = "SOPORTE"
)

// IsValid reports whether r is one of the known platform roles.
func (r SystemRole) IsValid() bool {
	switch r {
	case SystemRoleSuperAdmin, SystemRoleAdmin, SystemRoleSupport:
		return true
	}
	return false
}

func (SystemUser) TableName() string {
	return "system_users"
}

func (u SystemUser) FullName() string {
	return u.FirstName + " " + u.LastName
}
