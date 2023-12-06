package routes

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"strings"

	"github.com/gofiber/fiber/v2"
	controllers "tafakor.app/controllers"
	"tafakor.app/utils"
)

type publishmentBody struct {
	PostingType   string `query:"posting_type" json:"posting_type"`
	FileURL       string `query:"file_url" json:"file_url"`
	VerseID       string `query:"verse_id" json:"verse_id"`
	StockID       string `query:"stock_id" json:"stock_id"`
	StockProvider string `query:"stock_provider" json:"stock_provider"`
}

func PublishRoutes(group fiber.Router, db *sql.DB) {
	group.Post("/", func(c *fiber.Ctx) error {
		// Enviroment Variables
		SENDER_EMAIL := os.Getenv("SENDER_EMAIL")
		SENDER_PASS := os.Getenv("SENDER_PASS")
		EMAIL_HOST := os.Getenv("EMAIL_HOST")
		SMTP_PORT := os.Getenv("SMTP_PORT")
		SUPERVISOR_EMAIL := os.Getenv("SUPERVISOR_EMAIL")
		TAFAKOR_ENDPOINT := os.Getenv("TAFAKOR_ENDPOINT")

		var body publishmentBody
		json.Unmarshal(c.Body(), &body)

		parameters := fmt.Sprintf("?posting_type=%v&file_url=%v&verse_id=%v&stock_id=%v&stock_provider=%v", body.PostingType, body.FileURL, body.VerseID, body.StockID, body.StockProvider)

		acceptLink := fmt.Sprintf("%v/publish/accept", TAFAKOR_ENDPOINT) + parameters
		rejectLink := fmt.Sprintf("%v/publish/reject", TAFAKOR_ENDPOINT) + parameters
		rerenderLink := fmt.Sprintf("%v/publish/reject-stock", TAFAKOR_ENDPOINT) + parameters

		r := strings.NewReplacer("|POST-LINK|", body.FileURL, "|ACCEPT-LINK|", acceptLink, "|REJECT-LINK|", rejectLink, "|REJECT-RENDER-LINK|", rerenderLink)

		template, err := os.ReadFile("templates/approval.html")

		if err != nil {
			log.Fatal(err)
		}

		emailBody := r.Replace(string(template))

		err = utils.SendMail(SENDER_EMAIL, SENDER_PASS, EMAIL_HOST, SMTP_PORT, SUPERVISOR_EMAIL, "منشور جديد قيد الموافقة", emailBody)

		if err != nil {
			log.Fatal(err)
		}

		return c.JSON(nil)
	})

	group.Get("/accept", func(c *fiber.Ctx) error {
		var parameters publishmentBody
		c.QueryParser(&parameters)

		status, id := controllers.PublishToFB(parameters.PostingType, parameters.FileURL)

		if status == true {
			postQuery := "INSERT INTO post(verse, published, state, footage, postid) VALUES( $1, $2, $3, $4, $5 ) RETURNING id"
			postInsert, _ := db.Prepare(postQuery)

			var postId string
			postInsert.QueryRow(parameters.VerseID, true, "accepted", parameters.FileURL, id).Scan(postId)

			stockQuery := "INSERT INTO stock_footage(id, post, provider, state) VALUES( $1, $2, $3, $4 )"
			stockInsert, _ := db.Prepare(stockQuery)

			stockInsert.Exec(parameters.StockID, postId, parameters.StockProvider, "accepted")

			postInsert.Close()
			stockInsert.Close()

			return c.JSON(id)
		} else {
			return c.JSON(nil)
		}
	})

	group.Get("/reject", func(c *fiber.Ctx) error {
		var parameters publishmentBody
		c.QueryParser(&parameters)

		postQuery := "INSERT INTO post(verse, published, state, footage, postid) VALUES( $1, $2, $3, $4, $5 ) RETURNING id"
		postInsert, _ := db.Prepare(postQuery)

		var postId string
		postInsert.QueryRow(parameters.VerseID, true, "rejected", parameters.FileURL, nil).Scan(postId)

		stockQuery := "INSERT INTO stock_footage(id, post, provider, state) VALUES( $1, $2, $3, $4 )"
		stockInsert, _ := db.Prepare(stockQuery)

		stockInsert.Exec(parameters.StockID, postId, parameters.StockProvider, "rejected")

		postInsert.Close()
		stockInsert.Close()

		return c.JSON(true)
	})

	// group.Get("/reject-stock", func(c *fiber.Ctx) error {

	// })
}
