package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"tafakor.app/utils"
)

// Response returned by FB when after publishment
type account struct {
	ID          string `json:"id"`
	AccessToken string `json:"access_token"`
}

// Response returned by FB when after publishment
type accounts struct {
	Data []account `json:"data"`
}

// Response returned by FB when after publishment
type postPublishmentResponse struct {
	ID string `json:"id"`
}
type reelPublishmentResponse struct {
	Success bool   `json:"success"`
	ID      string `json:"post_id"`
}

type igReelPublishmentResponse struct {
	ID string `json:"id"`
}

type reelSessionInitRequest struct {
	UploadPhase string `json:"upload_phase"`
	AccessToken string `json:"access_token"`
}

type reelSessionInitResponse struct {
	VideoID string `json:"video_id"`
}

// var HASH_TAGS = "أَلا بِذِكرِ اللَّهِ تَطمَئِنُّ القُلوبُ 🤍🤍 \n #القران_الكريم #القرآن #قرآن #quran #تفكر #remembrance"
var HASH_TAGS = " "

/*
@desc Publishes Facaebook Posts
@param token - Temp FB access token
@param fileURL - Uploaded file url
*/
func FBPost(token string, fileURL string) (bool, string) {
	// Enviroment Variables
	var TAFAKOR_ID string = os.Getenv("TAFAKOR_ID")

	// Publishment Data
	data := url.Values{}
	data.Add("access_token", token)
	data.Add("file_url", fileURL)
	data.Add("message", HASH_TAGS)

	// Endpoint for post publishment
	videoPostEndpoint := fmt.Sprintf("https://graph-video.facebook.com/v18.0/%v/videos", TAFAKOR_ID)

	// Publishing post
	publishmentRes, _ := http.PostForm(videoPostEndpoint, data)
	publishmentStatus := utils.ParseJSONResponses[postPublishmentResponse](publishmentRes.Body)

	status := false
	if publishmentStatus.ID != "" {
		status = true
	}

	fmt.Println("Status:", status)

	// Publishment status
	return status, publishmentStatus.ID
}

/*
@desc Publishes Facaebook Reels
@param token - Temp FB access token
@param fileURL - Uploaded file url
*/
func FBReel(token string, fileURL string) (bool, string) {
	// Enviroment Variables
	var TAFAKOR_ID string = os.Getenv("TAFAKOR_ID")

	// Initting http client
	client := &http.Client{}

	// Session Init Data
	sessionInitEndpoint := fmt.Sprintf("https://graph.facebook.com/v18.0/%v/video_reels", TAFAKOR_ID)
	sessionInitBody, _ := json.Marshal(reelSessionInitRequest{UploadPhase: "start", AccessToken: token})
	// Initting Session
	sessionInitResponse, _ := http.Post(sessionInitEndpoint, "application/json", bytes.NewBuffer(sessionInitBody))

	// Initiated Video ID
	videoId := utils.ParseJSONResponses[reelSessionInitResponse](sessionInitResponse.Body).VideoID

	fmt.Println("VideoId", videoId)

	// Video Upload Endpoint
	videoUploadEndpoint := fmt.Sprintf("https://rupload.facebook.com/video-upload/v18.0/%v", videoId)
	fmt.Println("videoUploadEndpoint", videoUploadEndpoint)

	// Initiating Upload Request
	req, _ := http.NewRequest("POST", videoUploadEndpoint, nil)

	// Initiating Upload Data
	req.Header.Add("Authorization", "OAuth "+token)
	req.Header.Add("file_url", fileURL)
	req.Body = nil

	// Requesting Upload
	client.Do(req)

	// fmt.Println(utils.ParseJSONResponses[any](resp.Body))

	// Endpoint for post publishment
	reelPublishEndpoint := fmt.Sprintf("https://graph.facebook.com/v18.0/%v/video_reels?access_token=%v&video_id=%v&upload_phase=finish&video_state=PUBLISHED", TAFAKOR_ID, token, videoId)

	// Publishing post
	statusResponse, _ := http.Post(reelPublishEndpoint, "application/json", nil)
	respo := utils.ParseJSONResponses[reelPublishmentResponse](statusResponse.Body)

	// Publishment status
	return respo.Success, respo.ID
}

func _igReelPublish(uploadId string) string {
	var TAFAKOR_ID_INSTAGRAM string = os.Getenv("TAFAKOR_ID_INSTAGRAM")
	var USER_ACCESS_TOKEN string = os.Getenv("USER_ACCESS_TOKEN")

	reelPublishEndpoint := fmt.Sprintf("https://graph.facebook.com/v18.0/%v/media_publish?creation_id=%v&access_token=%v", TAFAKOR_ID_INSTAGRAM, uploadId, USER_ACCESS_TOKEN)

	fmt.Println(reelPublishEndpoint)
	nullBody := strings.NewReader("!")
	publishmentRes, _ := http.Post(reelPublishEndpoint, "", nullBody)
	fmt.Println(publishmentRes.StatusCode)

	// Printing Response
	// fmt.Println(publishmentRes.Body)

	fmt.Println(publishmentRes.Body.Read([]byte{}))
	if publishmentRes.StatusCode != 200 {
		return "NOT_YET"
	} else {
		PublishmentStatus := utils.ParseJSONResponses[igReelPublishmentResponse](publishmentRes.Body)
		return PublishmentStatus.ID
	}

}

func IGReel(fileURL string) string {
	var TAFAKOR_ID_INSTAGRAM string = os.Getenv("TAFAKOR_ID_INSTAGRAM")
	var USER_ACCESS_TOKEN string = os.Getenv("USER_ACCESS_TOKEN")

	// Endpoint for post publishment
	reelUploadEndpoint := fmt.Sprintf("https://graph.facebook.com/v18.0/%v/media?video_url=%v&access_token=%v&media_type=REELS&thumb_offset=2000&caption=%v", TAFAKOR_ID_INSTAGRAM, fileURL, USER_ACCESS_TOKEN, url.QueryEscape(HASH_TAGS))

	fmt.Println(reelUploadEndpoint)

	// Uploading Reel
	nullBody := strings.NewReader("!")
	uploadingRes, err := http.Post(reelUploadEndpoint, "", nullBody)

	if err == nil {
		uploadingStatus := utils.ParseJSONResponses[igReelPublishmentResponse](uploadingRes.Body)
		fmt.Println(uploadingStatus.ID + " instagaram")

		var publishmentId string

		publishmentId = _igReelPublish(uploadingStatus.ID)

		for publishmentId == "NOT_YET" {
			fmt.Println(publishmentId)
			time.Sleep(10 * time.Second)
			publishmentId = _igReelPublish(uploadingStatus.ID)
		}

		return publishmentId
	} else {
		log.Fatal(err)
	}

	return "ERROR"
}

/*
@desc Publishes Media to Social Media (reels | posts)
@param postURL - Uploaded post url
@param reelURL - Uploaded reel url
*/
func SocialPublishment(postURL string, reelURL string, platform string) (bool, string) {
	// Enviroment Variables
	var USER_ACCESS_TOKEN string = os.Getenv("USER_ACCESS_TOKEN")
	var TAFAKOR_ID string = os.Getenv("TAFAKOR_ID")
	var USER_ID string = os.Getenv("USER_ID")

	var postingTypes = [2]string{"reel", "post"}
	randomIndex := rand.Intn(2)

	facebookPublishmentType := postingTypes[randomIndex]

	// Accounts Request endpoint
	tokensEndpoint := fmt.Sprintf("https://graph.facebook.com/v18.0/%v/accounts?access_token=%v", USER_ID, USER_ACCESS_TOKEN)

	// Requesting Session Token
	tokensResponse, _ := http.Get(tokensEndpoint)
	var accounts accounts = utils.ParseJSONResponses[accounts](tokensResponse.Body) // Parsing Response

	// Saving Token
	var token string
	for _, account := range accounts.Data {
		if account.ID == TAFAKOR_ID {
			token = account.AccessToken
		}
	}

	var status bool = false
	var id string

	// Facebook Publish
	if platform == "all" || platform == "facebook" {
		fmt.Println("facebook")
		switch facebookPublishmentType {
		case "reel":
			fmt.Println("Reel")
			reelStatus, reelID := FBReel(token, reelURL)
			status = reelStatus
			id = reelID
		case "post":
			fmt.Println("Post")
			postStatus, postID := FBPost(token, postURL)
			status = postStatus
			id = postID
		}
	}

	if platform == "all" || platform == "instagram" {
		fmt.Println("Intsgaram")
		// Instagram Publish
		igStatus := IGReel(reelURL)

		if igStatus != "ERROR" {
			status = true
		}
	}

	return status, id
}
