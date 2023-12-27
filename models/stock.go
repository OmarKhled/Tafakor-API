package models

type Stock struct {
	ID       string `json:"id"`
	StockID  string `json:"stockid"`
	PostID   int    `json:"post"`
	Provider string `json:"provider"`
	State    string `json:"state"`
}
