package model

import "time"

type User struct {
	ID        uint      `gorm:"primaryKey;" json:"id"`
	Username  string    `gorm:"not null;" json:"username"`
	Email     string    `gorm:"unique;uniqueIndex;not null;" json:"email"`
	Password  string    `gorm:"" json:"password,omitempty"`
	Role      string    `gorm:"not null;" json:"role"`
	FullName  string    `json:"full_name"`
	Age       int       `json:"age"`
	ImageURL  string    `json:"image_url"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`

	UserSetting UserSetting `json:"user_setting,omitempty"`
}
