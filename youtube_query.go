package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

type YoutubeQuery struct {
	BaseURL    string
	APIKey     string
	ChannelID  string
	Parts      []string
	OrderBy    string
	MaxResults int
	PageToken  string
	Type       string
}

func NewYoutubeQuery(channelID string, pageToken string) *YoutubeQuery {

	r := new(YoutubeQuery)
	r.BaseURL = "https://www.googleapis.com/youtube/v3/search"
	r.APIKey = os.Getenv("YT_API_KEY")
	r.ChannelID = channelID
	r.Parts = append(r.Parts, "snippet", "id")
	r.OrderBy = "date"
	r.MaxResults = 50
	r.PageToken = pageToken
	r.Type = "video"
	return r
}

func (yq *YoutubeQuery) Run() *ChannelVideos {

	parts := strings.Join(yq.Parts, ",")
	url := fmt.Sprintf("%v?key=%v&channelId=%v&type=%v&part=%v&order=%v&maxResults=%v", yq.BaseURL, yq.APIKey, yq.ChannelID, yq.Type, parts, yq.OrderBy, yq.MaxResults)
	if yq.PageToken != "" {
		url += fmt.Sprintf("&pageToken=%v", yq.PageToken)
	}
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	var videos ChannelVideos
	json.Unmarshal(body, &videos)

	if err != nil {
		panic(err)
	}
	return &videos
}
