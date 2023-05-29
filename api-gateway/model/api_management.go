package model

import (
	"time"
)

type APIManagement struct {
	ID                uint      `gorm:"primaryKey;" json:"id"`
	APIName           string    `gorm:"unique;not null;" json:"api_name"`
	ServiceName       string    `gorm:"not null;" json:"service_name"`
	EndpointURL       string    `gorm:"not null;" json:"endpoint_url"`
	HashedEndpointURL string    `gorm:"unique;uniqueIndex;not null;" json:"hashed_endpoint_url"`
	IsAvailable       bool      `gorm:"not null;default:false" json:"is_available"`
	CreatedAt         time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt         time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

type ShortenReq struct {
	APIName     string `json:"api_name" binding:"required"`
	ServiceName string `json:"service_name" binding:"required"`
	EndpointURL string `json:"endpoint_url" binding:"required,url"`
	IsAvailable bool   `json:"is_available" binding:"-"`
}
