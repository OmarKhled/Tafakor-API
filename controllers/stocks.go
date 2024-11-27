package controllers

import (
	"database/sql"
	"log"

	"tafakor.app/models"
)

type StockFootage struct {
	ID           int    `json:"id"`
	StockID      string `json:"stockid"`
	ProviderName string `json:"provider"`
	State        string `json:"state"`
}

func GetStocks(db *sql.DB, verseId int) []StockFootage {
	var stocks []StockFootage

	// Selection Query
	var rows *sql.Rows
	rows, _ = db.Query(`
		(
			SELECT id, stockid, provider, state FROM stock_footage
			WHERE state NOT IN ('rejected-once', 'pending', 'discarded') 
		)
				UNION
		(
			SELECT id, stockid, provider, state FROM stock_footage
			WHERE state = 'rejected-once' AND post NOT IN (SELECT id FROM post WHERE verse = $1)
		);
	`, verseId)

	// Looping through returned rows
	for rows.Next() {
		// Initiate new footage
		var footage StockFootage

		// Save new footage
		err := rows.Scan(&footage.ID, &footage.StockID, &footage.ProviderName, &footage.State)

		if err != nil {
			log.Fatal(err)
		}

		// Appending new footage
		stocks = append(stocks, footage)
	}

	return stocks
}

func RecordStock(db *sql.DB, id string, postID int, provider string, state string) string {
	// Inserted Post ID
	var stockID string

	// Directly execute the query using QueryRow (without Prepare)
	stockQuery := "INSERT INTO stock_footage(stockid, post, provider, state) VALUES( $1, $2, $3, $4 ) RETURNING id"

	var err error

	if postID != 0 {
		// Insertion Query
		err = db.QueryRow(stockQuery, id, postID, provider, state).Scan(&stockID)
	} else {
		// Insertion Query
		err = db.QueryRow(stockQuery, id, nil, provider, state).Scan(&stockID)
	}

	if err != nil {
		log.Fatal(err)
	}

	return stockID
}

// func RecordStock(db *sql.DB, id string, postID int, provider string, state string) string {
// 	// Inserted Stock ID
// 	var stockID string

// 	// Insertion Query Prepare
// 	stockQuery := "INSERT INTO stock_footage(stockid, post, provider, state) VALUES( $1, $2, $3, $4 ) RETURNING id"
// 	stockInsert, err := db.Prepare(stockQuery)

// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	if postID != 0 {
// 		// Insertion Query
// 		err = stockInsert.QueryRow(id, postID, provider, state).Scan(&stockID)
// 	} else {
// 		// Insertion Query
// 		err = stockInsert.QueryRow(id, nil, provider, state).Scan(&stockID)
// 	}

// 	fmt.Println(err)
// 	fmt.Println(stockID)
// 	stockInsert.Close()

// 	return stockID
// }

func GetStockByPostID(db *sql.DB, postID int) models.Stock {
	var stock models.Stock

	// Insertion Query Prepare
	postQuery := "SELECT id, stockid, post, provider, state FROM stock_footage WHERE post = $1"

	row := db.QueryRow(postQuery, postID)
	err := row.Scan(&stock.ID, &stock.StockID, &stock.PostID, &stock.Provider, &stock.State)

	if err != nil {
		log.Fatal(err)
	}

	return stock
}

func ChangeStockStatus(db *sql.DB, stockID string, stockStatus string) bool {
	// Insertion Query Prepare
	postQuery := "UPDATE stock_footage SET state = $1 WHERE id = $2"

	_, err := db.Exec(postQuery, stockStatus, stockID)

	if err != nil {
		log.Fatal(err)
	}

	return true
}

func DeleteStock(db *sql.DB, stockID string) bool {
	// Insertion Query Prepare
	postQuery := "DELETE FROM stock_footage WHERE id = $2"

	_, err := db.Exec(postQuery, stockID)

	if err != nil {
		log.Fatal(err)
	}

	return true
}
