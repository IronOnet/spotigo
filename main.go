package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/schollz/progressbar/v3"
	"github.com/urfave/cli/v2"
)

func main(){
	app := &cli.App{
		Name: "spotigo",
		Usage: "Download songs from Youtube and Spotify",
		Action: func(c *cli.Context) error{
			song := c.Args().Get(0)
			artist := c.Args().Get(1)
			downloadSong(song, artist)
			return nil
		},
		ArgsUsage: "<song> <artist>",
	}

	err := app.Run(os.Args)
	if err != nil{
		log.Fatal(err)
	}
}


func downloadSong(song, artist string){
	query := fmt.Sprintf("%s %s official audio", song, artist)

	// Search for the song on Youtube
	resp, err := http.Get("https://www.youtube.com/results?search_query="+ query)
	if err != nil{
		log.Fatalf("Failed to search Youtube: %v", err)
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil{
		log.Fatalf("Failed to parse search results: %v", err)
	}

	// Find the first Youtube video link
	videoLink := doc.Find("a.yt-uix-tile-link").First()
	if videoLink.Length() == 0{
		log.Fatalf("Video not found for: %s by %s", song, artist)
	}

	videoURL, _ := videoLink.Attr("href")
	videoURL = "https://www.youtube.com" + videoURL

	// Download the audio stream
	resp, err = http.Get(videoURL)
	if err != nil{
		log.Fatalf("Failed to retrieve Youtube video details: %v", err)
	}

	defer resp.Body.Close()

	doc, err = goquery.NewDocumentFromReader(resp.Body)
	if err != nil{
		log.Fatalf("Failed to parse Youtube video details: %v", err)
	}

	audioURL, _ := doc.Find("script:contains('audioDetails')").First().Html()
	audioURL = strings.Split(strings.Split(audioURL, `"audioDetails":{"adapttiveFormats":[{"url":`))[0]

	// Download the audio file
	resp, err := http.Get(audioURL)
	if err != nil{
		log.Fatalf("Failed to download the audio file: %v", err)
	}

	defer resp.Body.Close()

	file, err := os.Create(fmt.Sprintf("%s - %s.mp3", song, artist))

	if err != nil{
		log.Fatalf("Failed to create audio file: %v", err)
	}

	defer file.Close()

	bar := progressbar.DefaultBytes(
		resp.ContentLength,
		"Downloading"
	)

	bar.RenderBlank()

	_, err = bar.Copy(resp.Body, file)
	if err != nil{
		log.Fatalf("Failed to write audio file: %v", err)
	}

	fmt.Println("Download Complete")
}
