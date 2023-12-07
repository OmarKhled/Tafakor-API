package controllers

import (
	"database/sql"
)

type StockFootage struct {
	ID           string `json:"id"`
	PostID       string `json:"post"`
	ProviderName string `json:"provider"`
	State        string `json:"state"`
}

func GetStocks(db *sql.DB, postID int) []StockFootage {
	var stocks []StockFootage

	// Selection Query
	var rows *sql.Rows
	rows, _ = db.Query(`
		(
			SELECT * FROM stock_footage
			WHERE state != 'rejected-once'
		)
				UNION
		(
			SELECT * FROM stock_footage
			WHERE state = 'rejected-once' AND post = $1
		);
	`, postID)

	// Looping through returned rows
	for rows.Next() {
		// Initiate new footage
		var footage StockFootage

		// Save new footage
		rows.Scan(&footage.ID, &footage.PostID, &footage.ProviderName, &footage.State, nil)

		// Appending new footage
		stocks = append(stocks, footage)
	}

	return stocks
}

func RecordStock(db *sql.DB, id string, postID string, provider string, state string) string {
	// Inserted Stock ID
	var stockID string

	// Insertion Query Prepare
	stockQuery := "INSERT INTO stock_footage(id, post, provider, state) VALUES( $1, $2, $3, $4 ) RETURNING id"
	stockInsert, _ := db.Prepare(stockQuery)

	// Insertion Query
	stockInsert.QueryRow(id, postID, provider, state).Scan(stockID)
	stockInsert.Close()

	return stockID
}