package controllers

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"tafakor.app/utils"
)

// Requests approval on publishments from supervisor
func RequestPostApproval(postId int, postURL string, reelURL string) {
	// Enviroment Variables
	SENDER_EMAIL := os.Getenv("SENDER_EMAIL")
	SENDER_PASS := os.Getenv("SENDER_PASS")
	EMAIL_HOST := os.Getenv("EMAIL_HOST")
	SMTP_PORT := os.Getenv("SMTP_PORT")
	SUPERVISOR_EMAIL := os.Getenv("SUPERVISOR_EMAIL")
	TAFAKOR_ENDPOINT := os.Getenv("TAFAKOR_ENDPOINT")

	// Converting Body to query
	parameters := fmt.Sprintf("?post_id=%v&platform=", postId)

	// URLS Reuired by approval email
	acceptLink := fmt.Sprintf("%v/publish/accept", TAFAKOR_ENDPOINT) + parameters + "all"                        // |ACCEPT|
	acceptFacebookLink := fmt.Sprintf("%v/publish/accept", TAFAKOR_ENDPOINT) + parameters + "facebook"           // |ACCEPT-FACEBOOK|
	acceptInstagramLink := fmt.Sprintf("%v/publish/accept", TAFAKOR_ENDPOINT) + parameters + "instagram"         // |ACCEPT-INSTAGRAM|
	rejectLink := fmt.Sprintf("%v/publish/reject", TAFAKOR_ENDPOINT) + parameters + "all"                        // |REJECT|
	rejectStockLink := fmt.Sprintf("%v/publish/reject/stock", TAFAKOR_ENDPOINT) + parameters + "all"             // |REJECT-STOCK|
	rejectVerseLink := fmt.Sprintf("%v/publish/reject/verse", TAFAKOR_ENDPOINT) + parameters + "all"             // |REJECT-VERSE|
	rejectStockForPostLink := fmt.Sprintf("%v/publish/reject/stock-post", TAFAKOR_ENDPOINT) + parameters + "all" // |REJECT-STOCK-ONCE|

	// Email template
	resp, _ := http.Get("https://tafkor.s3.eu-central-1.amazonaws.com/assets/approval.html")

	fmt.Println(resp)

	// Template parsing
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	template := buf.String()

	// Email template string replacer
	r := strings.NewReplacer("|POST-LINK|", postURL, "|REEL-LINK|", reelURL, "|ACCEPT|", acceptLink, "|REJECT|", rejectLink, "|REJECT-STOCK|", rejectStockLink, "|REJECT-VERSE|", rejectVerseLink, "|REJECT-STOCK-ONCE|", rejectStockForPostLink, "|ACCEPT-FACEBOOK|", acceptFacebookLink, "|ACCEPT-INSTAGRAM|", acceptInstagramLink)

	// Template Filling
	emailBody := r.Replace(string(template))

	// Email Send
	err := utils.SendMail(SENDER_EMAIL, SENDER_PASS, EMAIL_HOST, SMTP_PORT, SUPERVISOR_EMAIL, "منشور جديد قيد الموافقة", emailBody)

	if err != nil {
		log.Fatal(err)
	}
}
