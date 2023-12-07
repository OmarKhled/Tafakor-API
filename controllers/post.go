package controllers

import "database/sql"

func RecordPost(db *sql.DB, verseID string, published bool, state string, footageURL string, publishmentid string) int {
	// Inserted Post ID
	var postId int

	// Insertion Query Prepare
	postQuery := "INSERT INTO post(verse, published, state, footageUrl, publishmentid) VALUES( $1, $2, $3, $4, $5 ) RETURNING id"
	postInsert, _ := db.Prepare(postQuery)

	// Insertion Query
	postInsert.QueryRow(verseID, published, state, footageURL, publishmentid).Scan(&postId)
	postInsert.Close()

	return postId
}
