package verses

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
)

type response (struct {
	ID          string         `json:"id"`
	SurahNumber int            `json:"surah_number"`
	From        int            `json:"from"`
	To          int            `json:"to"`
	Count       int            `json:"count"`
	LastPublish sql.NullString `json:"last_publish"`
})

func GetVerses(db *sql.DB) func(c *fiber.Ctx) error {

	return func(c *fiber.Ctx) error {
		var verses []response

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
		defer rows.Close()

		for rows.Next() {

			var verse response
			if err := rows.Scan(&verse.ID, &verse.SurahNumber, &verse.From, &verse.To, &verse.Count, &verse.LastPublish); err != nil {
				log.Fatal(err)
				return c.JSON(verses)
			}
			fmt.Println("verse", verse)
			verses = append(verses, verse)
		}

		fmt.Println(verses)

		return c.JSON(verses)
	}

}
