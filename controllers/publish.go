package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"

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

type reelSessionInitRequest struct {
	UploadPhase string `json:"upload_phase"`
	AccessToken string `json:"access_token"`
}

type reelSessionInitResponse struct {
	VideoID string `json:"video_id"`
}

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

	// Endpoint for post publishment
	videoPostEndpoint := fmt.Sprintf("https://graph-video.facebook.com/v18.0/%v/videos", TAFAKOR_ID)

	fmt.Println(videoPostEndpoint)
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

func strings(respo any) {
	panic("unimplemented")
}

/*
@desc Publishes Media to Facaebook (reels | posts)
@param publishmentType - type of reel or post
@param fileURL - Uploaded file url
*/
func PublishToFB(publishmentType string, fileURL string) (bool, string) {
	// Enviroment Variables
	var USER_ACCESS_TOKEN string = os.Getenv("USER_ACCESS_TOKEN")
	var TAFAKOR_ID string = os.Getenv("TAFAKOR_ID")
	var USER_ID string = os.Getenv("USER_ID")

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

	fmt.Println("Token:", token)

	var status bool = false
	var id string

	switch publishmentType {
	case "reel":
		fmt.Println("Reel")
		reelStatus, reelID := FBReel(token, fileURL)
		status = reelStatus
		id = reelID
	case "post":
		fmt.Println("Post")
		postStatus, postID := FBPost(token, fileURL)
		status = postStatus
		id = postID
	}

	return status, id
}
