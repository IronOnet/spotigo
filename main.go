package main

import (
	"fmt"
	"log"
	//"net/http"
	"os"
	//"strings"

	//"github.com/PuerkitoBio/goquery"
	//"github.com/schollz/progressbar/v3"
	"github.com/urfave/cli/v2"
)

func main(){
	app := &cli.App{
		Name: "spotigo",
		Usage: "Download songs from Youtube and Spotify",
		Action: func(c *cli.Context) error{
			if c.Args().Get(0) != "" && c.Args().Get(1) != ""{
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
	if err != nil{
		log.Fatal(err)
	}
}


func downloadSong(song, artist string){
	fmt.Printf("Downloading song %s by %s", song, artist)
	fmt.Println("\nDownload Complete")
}
