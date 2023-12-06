package migrate

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"tafakor.app/utils"
)

type filter struct {
	ID            int    `json:"id"`
	SurahNumber   int    `json:"surah_number"`
	From          int    `json:"from"`
	To            int    `json:"to"`
	IndicatorText string `json:"indicator_text"`
}

type ReflectPost struct {
	ID      int      `json:"id"`
	Filters []filter `json:"filters"`
}

type VersesResponse struct {
	Posts []ReflectPost `json:"posts"`
}

/*
*

  - Tafakor's Database is mainly based on postings from https://quranreflect.com
  - This function is used to fetch all the postings used by them

*
*/

func FetchPosts() []ReflectPost {

	// All posts are saved here
	var posts []ReflectPost

	// Checking if local copy exists
	content, err := os.ReadFile("verses-data.json")

	if err != nil { // Local copy doesn't exist, fetching from API
		// Max number of fetched posts
		MAX_POSTS := 40

		// Counter for quranreflect api's pages pagenation
		page := 1

		// Looping tell all posts are fetched
		for len(posts) < MAX_POSTS {
			// Timeout between each consumed request rate
			timeout := 13

			// Holds the status for page fetch
			done := false

			// Loop until fetch done
			for !done {
				// API Endpoint with the pagination
				url := fmt.Sprintf("https://quranreflect.com/posts.json?client_auth_token=tUqQpl4f87wIGnLRLzG61dGYe03nkBQj&page=%v&tab=trending&lang=ar&featured=true", page)

				// Endpoint Get Request
				resp, err := http.Get(url)

				if err != nil { // Request Error
					log.Fatal(err)
				}

				// Closing request
				defer resp.Body.Close()

				// Logs
				fmt.Print("Page: ", page)
				println("  Status Code: ", resp.StatusCode)

				if resp.StatusCode == 200 { // If Success

					// Parsing JSON
					res := utils.ParseJSONResponses[VersesResponse](resp.Body)

					if len(res.Posts) != 0 {
						posts = append(posts, res.Posts...) // Saving fetched posts
						done = true                         // Mark page as done
						page++
					}
				} else { // If Request failed
					// Timeout Increse
					fmt.Printf("Timeout: %v seconds\n", timeout)
					time.Sleep(time.Duration(timeout) * time.Second)
					timeout += 1
				}

			}
		}

		// Saving a local copy of the data
		jsoncontent, _ := json.MarshalIndent(posts, "", "    ")
		err := os.WriteFile("verses-data.json", jsoncontent, 0644)

		if err != nil {
			log.Fatal(err)
		}
	} else { // Local Copy exists
		json.Unmarshal(content, &posts)
	}

	// Return all posts
	return posts
}
