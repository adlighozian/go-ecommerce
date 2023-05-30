package model

type RefreshToken struct {
	RefreshToken string `json:"refresh_token,omitempty"`
	UserID       uint   `json:"user_id,omitempty"`
	UserRole     string `json:"user_role,omitempty"`
}
