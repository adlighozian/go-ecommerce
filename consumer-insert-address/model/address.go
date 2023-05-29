package model

type Address struct {
	Id        int    `json:"id"`
	UserID    int    `json:"user_id"`
	Street    string `json:"street"`
	City      string `json:"city"`
	State     string `json:"state"`
	Country	  string `json:"country"`
	Zipcode   string `json:"zipcode"`
	PhoneNumber string `json:"phone_number"`
}

type AddressRequest struct {
	UserID    int    `json:"user_id"`
	Street    string `json:"street"`
	City      string `json:"city"`
	Country	  string `json:"country"`
	State     string `json:"state"`
	Zipcode   string `json:"zipcode"`
	PhoneNumber string `json:"phone_number"`
}