package migrate

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
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

type CitationText struct {
	ID   int    `json:"citation_id"`
	Text string `json:"text"`
}

type ReflectPost struct {
	ID            int                    `json:"id"`
	Filters       []filter               `json:"filters"`
	CitationTexts map[int][]CitationText `json:"citation_texts"`
}

type VersesResponse struct {
	Posts []ReflectPost `json:"posts"`
}

type Verse struct {
	ID          string
	SurahNumber int
	From        int
	To          int
	Length      int
}

/*
*

  - Tafakor's Database is mainly based on postings from https://quranreflect.com
  - This function is used to fetch all the postings used by them

*
*/

func FetchPosts() []Verse {

	// All verses are saved here
	var verses []Verse

	// Checking if local copy exists
	content, err := os.ReadFile("verses-data.json")

	if err != nil { // Local copy doesn't exist, fetching from API
		// Max number of fetched verses
		MAX_VERSES := 2000

		// Counter for quranreflect api's pages pagenation
		page := 1

		// Looping tell all verses are fetched
		for len(verses) < MAX_VERSES {
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
						for _, verse := range res.Posts {
							for filterIndex, meta := range verse.Filters {
								if len(verse.CitationTexts[filterIndex]) > 0 {
									verse := Verse{
										ID:          fmt.Sprintf("%v:%v", verse.ID, meta.ID),
										SurahNumber: meta.SurahNumber,
										From:        meta.From,
										To:          meta.To,
										Length:      len(strings.Fields(verse.CitationTexts[filterIndex][0].Text)),
									}
									// fmt.Printf("%+v\n", verse)

									verses = append(verses, verse)
								}

							}
						}
						done = true // Mark page as done
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
		jsoncontent, _ := json.MarshalIndent(verses, "", "    ")
		err := os.WriteFile("verses-data.json", jsoncontent, 0644)

		if err != nil {
			log.Fatal(err)
		}
	} else { // Local Copy exists
		json.Unmarshal(content, &verses)
	}

	// Return all posts
	return verses
}
