package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

func TriggerRender() {
	GITHUB_API_TOKEN := os.Getenv("GITHUB_API_TOKEN")

	client := &http.Client{}
	workFlowDispatchEndpoint := fmt.Sprintf("https://api.github.com/repos/omarkhled/tafakor/actions/workflows/render-video.yml/dispatches")

	reqBody := struct {
		Ref string `json:"ref"`
	}{Ref: "main"}
	var reqBodyBuffer bytes.Buffer
	json.NewEncoder(&reqBodyBuffer).Encode(reqBody)

	req, _ := http.NewRequest("POST", workFlowDispatchEndpoint, &reqBodyBuffer)

	req.Header.Add("Authorization", "Bearer "+GITHUB_API_TOKEN)
	req.Header.Add("X-GitHub-Api-Version", "2022-11-28")
	req.Header.Add("Accept", "application/vnd.github+json")
	req.Header.Add("Content-Type", "application/json")

	resp, _ := client.Do(req)

	content, _ := json.MarshalIndent(resp.Body, "", "")
	fmt.Println(string(content))
	fmt.Println(resp.StatusCode)

}
