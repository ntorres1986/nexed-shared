package models

import "time"

// RoleName enumerates the system (built-in) roles seeded at startup.
// Custom roles created later through the admin UI are plain data rows
// in `roles` and are not required to match one of these constants.
type RoleName string

const (
	RoleSuperAdmin     RoleName = "SUPER_ADMIN"
	RolePadre          RoleName = "PADRE"
	RoleDocente        RoleName = "DOCENTE"
	RoleAdministrativo RoleName = "ADMINISTRATIVO"
)

// Role groups a set of Permissions (RBAC). A User has one Role; a Role has
// many Permissions via the role_permissions join table.
//
// Roles are institution-scoped: each institution gets its own independent
// PADRE/DOCENTE/ADMINISTRATIVO rows (cloned with default permissions when
// the institution is created), so editing "Docente" for one school never
// affects any other school. Permission definitions themselves stay global
// (they name actual API capabilities, e.g. "payments.view") — only the
// role->permission assignment is per-institution.
// InstitutionID has no `not null` tag on purpose, even though every role
// is meant to belong to exactly one institution: this table pre-dates
// per-institution roles, and a populated `roles` table can't gain a new
// NOT NULL column outright (see database/legacy_migration.go, which adds
// it nullable and backfills it — the exact same approach already used for
// User.RoleID/InstitutionID). The guarantee is enforced in application
// code (RoleAdminService), not the schema.
type Role struct {
	ID            uint         `gorm:"primaryKey" json:"id"`
	InstitutionID uint         `gorm:"column:institution_id;uniqueIndex:idx_roles_institution_name;index" json:"institution_id"`
	Institution   *Institution `gorm:"foreignKey:InstitutionID" json:"-"`
	Name          string       `gorm:"column:name;size:50;not null;uniqueIndex:idx_roles_institution_name" json:"name"`
	IsSystem      bool         `gorm:"column:is_system;not null" json:"is_system"`
	Permissions   []Permission `gorm:"many2many:role_permissions;" json:"permissions,omitempty"`
	CreatedAt     time.Time    `gorm:"column:created_at" json:"created_at"`
	UpdatedAt     time.Time    `gorm:"column:updated_at" json:"updated_at"`
}

func (Role) TableName() string {
	return "roles"
}

// PermissionNames returns the flat list of permission name strings granted
// by this role. Role.Permissions must be preloaded before calling this.
func (r Role) PermissionNames() []string {
	names := make([]string, 0, len(r.Permissions))
	for _, p := range r.Permissions {
		names = append(names, p.Name)
	}
	return names
}
