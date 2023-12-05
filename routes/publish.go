package routes

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/gofiber/fiber/v2"
	controllers "tafakor.app/controllers"
)

type publishmentBody struct {
	PostingType   string `json:"posting_type"`
	FileURL       string `json:"file_url"`
	VerseID       string `json:"verse_id"`
	StockID       string `json:"stock_id"`
	StockProvider string `json:"stock_provider"`
}

func PublishRoutes(group fiber.Router, db *sql.DB) {
	group.Post("/", func(c *fiber.Ctx) error {
		var body publishmentBody
		json.Unmarshal(c.Body(), &body)

		status, id := controllers.PublishToFB(body.PostingType, body.FileURL)

		fmt.Println(status, id)

		if status == true {
			fmt.Println("tany")
			query := "INSERT INTO post(verse, published, footage, footageid) VALUES( $1, $2, $3, $4 )"
			insert, err := db.Prepare(query)

			if err != nil {
				fmt.Println(err)
			}

			resp, err := insert.Exec(body.VerseID, true, body.FileURL, string(id))
			fmt.Println(resp, err)

			insert.Close()
		}

		return c.JSON(status)
	})
}
