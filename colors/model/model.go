package model

import (
	"time"
)

type Colors struct {
	Id         int       `json:"id"`
	Name       string    `json:"name"`
	Created_at time.Time `json:"created_at"`
	Update_at  time.Time `json:"updated_at"`
}

type ColorsReq struct {
	Name string `json:"name"`
}

type Respon struct {
	Status int
	Data   any
}

type ResponSuccess struct {
	Status  int
	Message string
	Data    any
}
type ResponError struct {
	Status  int
	Message string
	Error   string
}
