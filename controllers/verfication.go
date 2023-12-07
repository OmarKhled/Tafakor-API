package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"tafakor.app/models"
	"tafakor.app/utils"
)

// Requests approval on publishments from supervisor
func RequestPostApproval(bodyBytes []byte) {
	// Enviroment Variables
	SENDER_EMAIL := os.Getenv("SENDER_EMAIL")
	SENDER_PASS := os.Getenv("SENDER_PASS")
	EMAIL_HOST := os.Getenv("EMAIL_HOST")
	SMTP_PORT := os.Getenv("SMTP_PORT")
	SUPERVISOR_EMAIL := os.Getenv("SUPERVISOR_EMAIL")
	TAFAKOR_ENDPOINT := os.Getenv("TAFAKOR_ENDPOINT")

	// Extracting Body content
	var body models.PublishmentParamaters
	json.Unmarshal(bodyBytes, &body)

	// Converting Body to query
	parameters := fmt.Sprintf("?posting_type=%v&file_url=%v&verse_id=%v&stock_id=%v&stock_provider=%v", body.PostingType, body.FileURL, body.VerseID, body.StockID, body.StockProvider)

	// URLS Reuired by approval email
	acceptLink := fmt.Sprintf("%v/publish/accept", TAFAKOR_ENDPOINT) + parameters
	rejectLink := fmt.Sprintf("%v/publish/reject", TAFAKOR_ENDPOINT) + parameters                        // |REJECT|
	rejectStockLink := fmt.Sprintf("%v/publish/reject/stock", TAFAKOR_ENDPOINT) + parameters             // |REJECT-STOCK|
	rejectVerseLink := fmt.Sprintf("%v/publish/reject/verse", TAFAKOR_ENDPOINT) + parameters             // |REJECT-VERSE|
	rejectStockForPostLink := fmt.Sprintf("%v/publish/reject/stock-post", TAFAKOR_ENDPOINT) + parameters // |REJECT-STOCK-ONCE|

	// Email template
	resp, _ := http.Get("https://tafakor.s3.eu-north-1.amazonaws.com/assets/approval.html")

	// Template parsing
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	template := buf.String()

	// Email template replacer
	r := strings.NewReplacer("|POST-LINK|", body.FileURL, "|ACCEPT|", acceptLink, "|REJECT|", rejectLink, "|REJECT-RENDER-LINK|", rejectStockLink, "|REJECT-VERSE|", rejectVerseLink, "|REJECT-STOCK-ONCE|", rejectStockForPostLink)

	// Template Filling
	emailBody := r.Replace(string(template))

	// Email Send
	err := utils.SendMail(SENDER_EMAIL, SENDER_PASS, EMAIL_HOST, SMTP_PORT, SUPERVISOR_EMAIL, "منشور جديد قيد الموافقة", emailBody)

	if err != nil {
		log.Fatal(err)
	}
}
