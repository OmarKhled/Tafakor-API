package models

type Post struct {
	ID            int    `json:"id"`
	VerseID       string `json:"verse"`
	Published     bool   `json:"published"`
	State         string `json:"state"`
	PublishmentID string `json:"publishmentid"`
	PostURL       string `json:"posturl"`
	ReelURL       string `json:"reelurl"`
}
