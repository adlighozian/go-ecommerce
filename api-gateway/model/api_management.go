package model

import (
	"time"
)

type APIManagement struct {
	ID                uint      `gorm:"primaryKey" json:"id"`
	APIName           string    `gorm:"unique" json:"api_name"`
	ServiceName       string    `json:"service_name"`
	EndpointURL       string    `gorm:"type:text" json:"endpoint_url"`
	HashedEndpointURL string    `gorm:"type:varchar(255);unique_index" json:"hashed_endpoint_url"`
	IsAvailable       bool      `gorm:"default:false" json:"is_available"`
	CreatedAt         time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt         time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

type ShortenReq struct {
	APIName     string `json:"api_name" binding:"required"`
	ServiceName string `json:"service_name" binding:"required"`
	EndpointURL string `json:"endpoint_url" binding:"required,url"`
	IsAvailable bool   `json:"is_available" binding:"-"`
}
