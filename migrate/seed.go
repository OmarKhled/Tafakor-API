package migrate

import (
	"database/sql"
	"fmt"
	// "fmt"
)

func Seed(db *sql.DB) {
	var posts []ReflectPost = FetchPosts()

	for index, post := range posts {
		if len(post.Filters) > 0 {
			for _, meta := range post.Filters {
				query := "INSERT INTO verse(id, surah_number, start, \"end\") VALUES( $1, $2, $3, $4 )"
				// query := "SELECT * FROM verse"

				insert, err := db.Prepare(query)
				if err != nil {
					fmt.Println(err)
				}

				id := fmt.Sprintf("%v:%v", post.ID, meta.ID)
				resp, err := insert.Exec(id, meta.SurahNumber, meta.From, meta.To)
				fmt.Println("Index:", index, " Inserted", id)
				insert.Close()

				if err != nil {
					fmt.Println(err, resp)
				}

			}
		}
	}
}
