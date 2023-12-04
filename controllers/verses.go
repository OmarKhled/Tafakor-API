package verses

import (
	"database/sql"
	"log"
	"math/rand"

	"github.com/gofiber/fiber/v2"
)

type verse struct {
	ID          string `json:"id"`
	SurahNumber int    `json:"surah_number"`
	From        int    `json:"from"`
	To          int    `json:"to"`
	Count       int    `json:"count"`
}

func GetVerse(db *sql.DB) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		random := c.QueryBool("random")
		// Stores DB fetched verses
		var verses []verse

		verseIndex := 0

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

		if random {
			verseIndex = rand.Intn(len(verses))
		}

		// Final reaponse
		return c.JSON(verses[verseIndex])
	}
}
