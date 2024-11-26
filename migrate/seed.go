package migrate

import (
	"database/sql"
	"fmt"
)

func Seed(db *sql.DB) {
	// Fetch posts from Quran Reflect's API
	var verses []Verse = FetchPosts()

	if len(verses) > 0 {
		// Loop through posts
		for index, verse := range verses {
			if verse.Length > 7 && verse.SurahNumber >= 2 && verse.SurahNumber <= 78 {
				query := "INSERT INTO verse(id, surah_number, start, \"end\") VALUES( $1, $2, $3, $4 )"

				insert, err := db.Prepare(query)

				if err != nil {
					fmt.Println(err)
				}

				resp, err := insert.Exec(verse.ID, verse.SurahNumber, verse.From, verse.To)

				fmt.Println("Index:", index, " Inserted", verse.ID)

				// End instruction
				insert.Close()

				if err != nil {
					fmt.Println(err, resp)
				}
			} else {
				fmt.Println("Skipped Index:", index, "Length:", verse.Length, "To:", verse.To, "From:", verse.From)
			}
		}
	}

}
