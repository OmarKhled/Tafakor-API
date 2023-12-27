package models

type Stock struct {
	ID       string `json:"id"`
	PostID   int    `json:"post"`
	Provider string `json:"provider"`
	State    string `json:"state"`
}
