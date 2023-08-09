package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	//"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/schollz/progressbar/v3"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "spotigo",
		Usage: "Download songs from Youtube and Spotify",
		Action: func(c *cli.Context) error {
			if c.Args().Get(0) != "" && c.Args().Get(1) != "" {
				song := c.Args().Get(0)
				artist := c.Args().Get(1)
				downloadSong(song, artist)
				return nil
			}
			return nil
		},
		ArgsUsage: "<song> <artist>",
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func downloadSong(song, artist string) {
	fmt.Printf("\nDownloading song %s by %s", song, artist)

	// The algorithm
	// Take two string inputs from the Arg
	// Pass those string inputs to a search query
	// the search query returns a list of results
	// select the first video returned
	// download the audio stream from the video
	// selected
	// save the audio file in a folder specified in the --path argument

	query := fmt.Sprintf("%s+%s+official+audio", song, artist)

	// Search for the song on Youtube
	resp, err := http.Get("https://www.youtube.com/results?search_query=" + query)
	if err != nil {
		log.Fatalf("Failed to search Youtube: %v", err)
	}

	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body) 
	if err != nil{
		log.Fatalf("\nfailed to parse search results: %v", err)
	}

	// Find the first Youtube video Link 
	videoLink := doc.Find("a.yt-uix-tile-link").First() 
	if videoLink.Length() == 0{
		log.Fatalf("\nvideo not found for: %s by %s", song, artist)
	} else{
		fmt.Println("\nVideo found")
	}

	videoUrl, _ := videoLink.Attr("href") 
	videoUrl = "https://www.youtube.com" + videoUrl 

	// Download the audio stream 
	resp, err = http.Get(videoUrl) 
	if err != nil{
		log.Fatalf("\nFailed to retrieve the youtube video details: %v", err)
	}

	defer resp.Body.Close()

	file, err := os.Create(fmt.Sprintf("%s - %s.mp3", song, artist))
	if err != nil{
		log.Fatalf("Failed to create audio file: %v", err)
	}
	defer file.Close()

	bar := progressbar.DefaultBytes(
		resp.ContentLength, 
		"Downloading",
	)

	bar.RenderBlank()

	_, err = io.Copy(io.MultiWriter(file, bar), resp.Body)
	if err != nil{
		log.Fatalf("failed to write audio file: %v", err)
	}



	fmt.Println(string(bodyBytes))
	fmt.Println("\nDownload Complete")
}
