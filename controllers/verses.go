package versesController

import (
	"database/sql"
	"log"
	"math/rand"
)

type verse struct {
	ID          string `json:"id"`
	SurahNumber int    `json:"surah_number"`
	From        int    `json:"from"`
	To          int    `json:"to"`
	Count       int    `json:"count"`
}

func GetVerses(db *sql.DB) []verse {
	var verses []verse

	// Select query to get verses by priority
	rows, err := db.Query(`
		SELECT id, surah_number, start, "end",
			(SELECT count(id) FROM post WHERE V.id = post.verse) as count,
			(SELECT max(created_at) FROM post WHERE V.id = post.verse) as last_publish
		FROM verse as V
		ORDER BY count, last_publish DESC, created_at DESC;
	`)

	if err != nil {
		log.Fatal(err)
	}

	// Looping through returned rows
	for rows.Next() {
		// Initiate new verse
		var verse verse
		// Save new verse
		rows.Scan(&verse.ID, &verse.SurahNumber, &verse.From, &verse.To, &verse.Count, nil)
		// Appending new verse
		verses = append(verses, verse)
	}

	// Closing query
	defer rows.Close()

	return verses
}

func GetVerse(random bool, db *sql.DB) verse {
	// The returned verse indes (by default the first verse returned but can be changes to random verse based on the random property)
	verseIndex := 0

	// Fetch Verses
	verses := GetVerses(db)

	// Checks if random verse is requied
	if random {
		verseIndex = rand.Intn(len(verses))
	}

	return verses[verseIndex]
}
