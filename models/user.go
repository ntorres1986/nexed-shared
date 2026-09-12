package models

import "time"

// User represents an authenticated system user.
//
// InstitutionID is nil only for global (SUPER_ADMIN) users; every
// institution-bound user (PADRE/DOCENTE/ADMINISTRATIVO) must have it set.
// Authorization never trusts a client-supplied institution id — it is
// always resolved from this field via the authenticated session.
type User struct {
	ID            uint         `gorm:"primaryKey" json:"id"`
	InstitutionID *uint        `gorm:"column:institution_id;uniqueIndex:idx_users_institution_email;index" json:"institution_id"`
	Institution   *Institution `gorm:"foreignKey:InstitutionID" json:"institution,omitempty"`

	FirstName string `gorm:"column:first_name;size:100;not null" json:"first_name"`
	LastName  string `gorm:"column:last_name;size:100;not null" json:"last_name"`
	// Unique together with InstitutionID (see uniqueIndex tag above): the
	// same email may exist in two different institutions, but not twice
	// within the same one.
	Email    string `gorm:"column:email;size:150;not null;uniqueIndex:idx_users_institution_email" json:"email"`
	Password string `gorm:"column:password;size:255;not null" json:"-"`

	RoleID uint `gorm:"column:role_id;index" json:"role_id"`
	Role   Role `gorm:"foreignKey:RoleID" json:"role,omitempty"`

	IsActive bool `gorm:"column:is_active;not null" json:"is_active"`
	// Set for accounts provisioned by an admin (manual creation or CSV
	// import) with a temporary password. The mechanism (field + email
	// link) is in place; enforcing a forced change screen is left for a
	// follow-up phase, per project scope.
	MustChangePassword bool `gorm:"column:must_change_password;not null" json:"must_change_password"`

	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (User) TableName() string {
	return "users"
}

func (u User) FullName() string {
	return u.FirstName + " " + u.LastName
}

// IsSuperAdmin reports whether this user is the global platform admin
// (not bound to any institution).
func (u User) IsSuperAdmin() bool {
	return u.InstitutionID == nil
}
