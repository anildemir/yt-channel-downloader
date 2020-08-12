package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	ytdl "github.com/kkdai/youtube/downloader"
)

func main() {

	if len(os.Getenv("YT_API_KEY")) <= 0 {
		log.Fatal("Missing API key.")
	}

	channelID := flag.String("channel", "", "ID of the channel you would like to download.")

	flag.Parse()

	if len(*channelID) <= 0 {
		log.Fatalf("Missing argument, please enter a channel ID")
	}

	home := os.Getenv("YT_OUTPUT_DIR")

	if len(home) <= 0 {
		home, _ = os.UserHomeDir()
	}

	downloader := NewDownloader()

	var pageToken string

	for {
		query := NewYoutubeQuery(*channelID, pageToken)
		channelVideos := query.Run()

		channelTitle := channelVideos.Videos[0].Snippet.ChannelTitle
		totalResults := channelVideos.PageInfo.TotalResults

		videoIndex := 1
		downloader.OutputDir = filepath.Join(home, channelTitle)

		fmt.Println("Found", totalResults, "videos on channel:", channelTitle)
		fmt.Println("Starting to download:")

		for _, v := range channelVideos.Videos {

			videoName := fmt.Sprintf("%v.mp4", v.Snippet.Title)
			videoID := v.ID.VideoID

			fmt.Println(videoIndex, "/", totalResults, "|", time.Now(), "Downloading", videoName, "by", channelTitle)

			video, err := downloader.GetVideo(videoID)
			if err != nil {
				panic(err)
			}

			ffmpegVersionCmd := exec.Command("ffmpeg", "-version")
			if err := ffmpegVersionCmd.Run(); err != nil {
				log.Fatal(fmt.Errorf("Please check ffmpeg is installed correctly, err: %w", err))
			}

			if err := downloader.DownloadWithHighQuality(context.Background(), videoName, video, "hd1080"); err != nil {

				if err.Error() == "no Stream video/mp4 for hd1080 found" {
					fmt.Println("1080p version is not found for", videoName, " - Attempting 720p...")
					stream_index := -1
					for k, v := range video.Streams {
						if v.Quality == "hd720" && strings.HasPrefix(v.MimeType, "video/mp4") {
							stream_index = k
							break
						}
					}

					if stream_index >= 0 {
						fmt.Println("Found 720p, downloading...")
						if err := downloader.Download(context.Background(), video, &video.Streams[stream_index], videoName); err != nil {
							panic(err)
						}

					} else {
						fmt.Println("No quality found >=720p, skipping video.")
						continue
					}

				} else {
					panic(err)
				}
			}
			videoIndex++
		}

		if channelVideos.NextPageToken == "" {
			fmt.Println("All videos are downloaded.")
			break
		} else {
			pageToken = channelVideos.NextPageToken
		}

	}
}

func NewDownloader() ytdl.Downloader {
	downloader := ytdl.Downloader{}

	httpTransport := &http.Transport{
		IdleConnTimeout:       60 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
	httpTransport.DialContext = (&net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	}).DialContext
	downloader.HTTPClient = &http.Client{Transport: httpTransport}

	return downloader
}
