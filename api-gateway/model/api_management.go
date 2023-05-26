package model

import (
	"time"
)

type APIManagement struct {
	ID                uint   `gorm:"primaryKey"`
	APIName           string `gorm:"unique"`
	ServiceName       string
	EndpointURL       string    `gorm:"type:text"`
	HashedEndpointURL string    `gorm:"type:varchar(255);unique_index"`
	IsAvailable       bool      `gorm:"default:false"`
	CreatedAt         time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt         time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}

type ShortenReq struct {
	APIName     string `json:"api_name" binding:"required"`
	ServiceName string `json:"service_name" binding:"required"`
	EndpointURL string `json:"endpoint_url" binding:"required,url"`
	IsAvailable bool   `json:"is_available" binding:"-"`
}
