package model

type Wishlist struct {
	Id        int `json:"id"`
	UserID    int `json:"user_id"`
	ProductID int `json:"product_id"`
}

type WishlistRequest struct {
	UserID    int `json:"user_id"`
	ProductID int `json:"product_id"`
}