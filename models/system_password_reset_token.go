package models

import "time"

// SystemPasswordResetToken is the system_users equivalent of
// PasswordResetToken — kept as a separate table (not a shared/polymorphic
// one) so the system_users track has zero foreign-key coupling to the
// institution `users` table, matching the same separation principle as
// SystemUser itself.
type SystemPasswordResetToken struct {
	ID           uint       `gorm:"primaryKey" json:"id"`
	SystemUserID uint       `gorm:"column:system_user_id;not null;index" json:"system_user_id"`
	SystemUser   SystemUser `gorm:"foreignKey:SystemUserID;constraint:OnDelete:CASCADE" json:"-"`
	TokenHash    string     `gorm:"column:token_hash;size:255;not null;uniqueIndex" json:"-"`
	ExpiresAt    time.Time  `gorm:"column:expires_at;not null;index" json:"expires_at"`
	UsedAt       *time.Time `gorm:"column:used_at" json:"used_at"`
	CreatedAt    time.Time  `gorm:"column:created_at" json:"created_at"`
}

func (SystemPasswordResetToken) TableName() string {
	return "system_password_reset_tokens"
}

func (t SystemPasswordResetToken) IsExpired() bool {
	return time.Now().After(t.ExpiresAt)
}

func (t SystemPasswordResetToken) IsUsed() bool {
	return t.UsedAt != nil
}
