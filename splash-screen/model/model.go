package model

type SplashScreen struct {
	Id       int    `json:"id"`
	ImageURL string `json:"image_url"`
}

type SplashScreenRequest struct {
	ImageURL string `json:"image_url"`
}
