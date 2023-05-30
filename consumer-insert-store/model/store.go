package model

type Store struct {
	Id          int    `json:"id"`
	AddressID   int    `json:"address_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	ImageURL    string `json:"image_url"`
}

type StoreRequest struct {
	AddressID   int    `json:"address_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	ImageURL    string `json:"image_url"`
}
