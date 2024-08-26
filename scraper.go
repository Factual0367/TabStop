package internal

import (
	"encoding/json"
	"fmt"
	"strings"

	colly "github.com/gocolly/colly"
)

type Song struct {
	SongTitle    string `json:"title"`
	Artist       string `json:"artist"`
	ID           string `json:"id"`
	DownloadLink string `json:"downloadLink"`
}

func AddGPDownloadLinks(songs []Song) {
	for i := range songs {
		song := &songs[i]
		// new collector for each song to avoid adding same
		// link for all songs
		c := colly.NewCollector()

		c.OnResponse(func(r *colly.Response) {
			var data []map[string]interface{}
			err := json.Unmarshal(r.Body, &data)
			if err != nil {
				fmt.Println(err)
				return
			}

			if len(data) > 0 {
				if source, ok := data[0]["source"].(string); ok {
					song.DownloadLink = source
				}
			}
		})

		c.Visit("https://www.songsterr.com/api/meta/" + song.ID + "/revisions")
	}
}

func GetSongList(query string) []Song {
	headers := map[string]string{
		"User-Agent":                "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:109.0) Gecko/20100101 Firefox/119.0",
		"Accept":                    "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8",
		"Accept-Language":           "en-US,en;q=0.5",
		"DNT":                       "1",
		"Upgrade-Insecure-Requests": "1",
		"Sec-Fetch-Dest":            "document",
		"Sec-Fetch-Mode":            "navigate",
		"Sec-Fetch-Site":            "cross-site",
		"Connection":                "keep-alive",
	}

	c := colly.NewCollector()

	c.OnRequest(func(r *colly.Request) {
		for key, value := range headers {
			r.Headers.Set(key, value)
		}
	})

	var songs []Song

	pattern := query

	searchUrl := fmt.Sprintf("https://www.songsterr.com/?pattern=%s", strings.ReplaceAll(pattern, " ", "%20"))
	c.OnHTML("div[data-list='songs']", func(e *colly.HTMLElement) {
		e.ForEach("a", func(_ int, el *colly.HTMLElement) {
			var songTitle, artist string

			el.ForEach("div.B0c2e8", func(_ int, div *colly.HTMLElement) {
				songTitle = div.Text
			})

			el.ForEach("div.B0c21e", func(_ int, div *colly.HTMLElement) {
				artist = div.Text
			})

			href := el.Attr("href")
			id := strings.Split(href, "-s")[len(strings.Split(href, "-s"))-1]

			song := Song{songTitle, artist, id, id}
			songs = append(songs, song)
		})
	})

	c.Visit(searchUrl)
	return songs
}
