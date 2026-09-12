package models

import "time"

// PasswordResetToken represents a hashed, single-use password recovery token.
type PasswordResetToken struct {
	ID        uint       `gorm:"primaryKey" json:"id"`
	UserID    uint       `gorm:"column:user_id;not null;index" json:"user_id"`
	User      User       `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"-"`
	TokenHash string     `gorm:"column:token_hash;size:255;not null;uniqueIndex" json:"-"`
	ExpiresAt time.Time  `gorm:"column:expires_at;not null;index" json:"expires_at"`
	UsedAt    *time.Time `gorm:"column:used_at" json:"used_at"`
	CreatedAt time.Time  `gorm:"column:created_at" json:"created_at"`
}

// TableName overrides the default pluralized table name.
func (PasswordResetToken) TableName() string {
	return "password_reset_tokens"
}

// IsExpired reports whether the token has passed its expiration time.
func (t PasswordResetToken) IsExpired() bool {
	return time.Now().After(t.ExpiresAt)
}

// IsUsed reports whether the token has already been consumed.
func (t PasswordResetToken) IsUsed() bool {
	return t.UsedAt != nil
}
