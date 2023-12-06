package models

type PublishmentParamaters struct {
	PostingType   string `query:"posting_type" json:"posting_type"`
	FileURL       string `query:"file_url" json:"file_url"`
	VerseID       string `query:"verse_id" json:"verse_id"`
	StockID       string `query:"stock_id" json:"stock_id"`
	StockProvider string `query:"stock_provider" json:"stock_provider"`
}
