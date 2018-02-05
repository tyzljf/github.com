package models

type Job struct {
	ID int `json:"id"`
	Title string `json:"title"`
	Description string `json:"description"`
	UserID string `json:"user_id"`
}
