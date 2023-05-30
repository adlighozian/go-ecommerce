package model

import "time"

type UserSetting struct {
	ID           uint      `gorm:"primaryKey;" json:"id"`
	UserID       uint      `gorm:"not null;" json:"user_id"`
	Notification *bool     `gorm:"not null;default:true" json:"notification"`
	DarkMode     *bool     `gorm:"not null;default:false" json:"dark_mode"`
	LanguageID   uint      `gorm:"not null;default:1" json:"language_id"`
	CreatedAt    time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt    time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`

	Language Language `json:"language,omitempty"`
}
