package controllers

import (
	"database/sql"
	// "fmt"
	"log"

	"tafakor.app/models"
)

func RecordPost(db *sql.DB, verseID string, published bool, state string, postURL string, publishmentid int, reelURL string) int {
	// Inserted Post ID
	var postId int = 0

	// Directly execute the query using QueryRow (without Prepare)
	postQuery := "INSERT INTO post(verse, published, state, posturl, publishmentid, reelurl) VALUES( $1, $2, $3, $4, $5, $6 ) RETURNING id"
	row := db.QueryRow(postQuery, verseID, published, state, postURL, publishmentid, reelURL)

	// Scan the result into postId
	err := row.Scan(&postId)
	if err != nil {
		log.Fatal("Error executing the query: ", err)
	}

	return postId
}

func GetPost(db *sql.DB, postID int) models.Post {
	var post models.Post

	// Insertion Query Prepare
	postQuery := "SELECT id, verse, published, state, publishmentid, posturl, reelurl FROM post WHERE id = $1"

	row := db.QueryRow(postQuery, postID)
	err := row.Scan(&post.ID, &post.VerseID, &post.Published, &post.State, &post.PublishmentID, &post.PostURL, &post.ReelURL)

	if err != nil {
		log.Fatal(err)
	}

	return post
}

func ChangePostStatus(db *sql.DB, postID int, postStatus string) bool {
	// Insertion Query Prepare
	postQuery := "UPDATE post SET state = $1 WHERE id = $2"

	_, err := db.Exec(postQuery, postStatus, postID)

	if err != nil {
		log.Fatal(err)
	}

	return true
}

func DeletePost(db *sql.DB, postID int) bool {
	// Insertion Query Prepare
	postQuery := "DELETE FROM post WHERE id = $1"

	_, err := db.Exec(postQuery, postID)

	if err != nil {
		log.Fatal(err)
	}

	return true
}
