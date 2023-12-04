package migrate

import (
	"database/sql"
	"fmt"
)

func Seed(db *sql.DB) {
	// Fetch posts from Quran Reflect's API
	var posts []ReflectPost = FetchPosts()

	// Loop through posts
	for index, post := range posts {
		if len(post.Filters) > 0 {
			// For each segment in the post
			for _, meta := range post.Filters {
				// DB Insert Query
				query := "INSERT INTO verse(id, surah_number, start, \"end\") VALUES( $1, $2, $3, $4 )"

				insert, err := db.Prepare(query)

				if err != nil {
					fmt.Println(err)
				}

				// Final ID is combination of the segment ID and the post ID
				id := fmt.Sprintf("%v:%v", post.ID, meta.ID)
				// Insert into DB
				resp, err := insert.Exec(id, meta.SurahNumber, meta.From, meta.To)
				fmt.Println("Index:", index, " Inserted", id)

				// End instruction
				insert.Close()

				if err != nil {
					fmt.Println(err, resp)
				}

			}
		}
	}
}
