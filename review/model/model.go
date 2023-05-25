package model

type Review struct {
	Id        	int `json:"id"`
	UserID    	int `json:"user_id"`
	ProductID 	int `json:"product_id"`
	Rating		int `json:"rating_id"`
	ReviewText	int `json:"review_text"`
}

type ReviewRequest struct {
	UserID    	int `json:"user_id"`
	ProductID 	int `json:"product_id"`
	Rating		int `json:"rating_id"`
	ReviewText	int `json:"review_text"`
}