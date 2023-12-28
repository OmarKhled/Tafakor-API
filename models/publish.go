package models

type PublishmentParamaters struct {
	PostURL       string `query:"post_url" json:"post_url"`
	ReelURL       string `query:"reel_url" json:"reel_url"`
	VerseID       string `query:"verse_id" json:"verse_id"`
	StockID       string `query:"stock_id" json:"stock_id"`
	StockProvider string `query:"stock_provider" json:"stock_provider"`
}

type EmailSubmissionParameters struct {
	PostID   int    `query:"post_id" json:"post_id"`
	Platform string `query:"platform" json:"platform"`
}
