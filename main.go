package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"fyne.io/fyne/theme"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/gin-gonic/gin"
	"github.com/onurhanak/songsterrapi"
)

type Tab struct {
	id           string
	downloadLink string
	title        string
	artist       string
}

func startServer() {
	router := gin.Default()
	router.GET("/search", songsterrapi.SearchRequest)
	router.Run("localhost:8080")
}

func startServerInBackground() {
	go startServer()

}

func getTabs(query string) []Tab {

	resp, err := http.Get("http://localhost:8080/search?query=" + query)
	if err != nil {
		fmt.Println(err)

	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)

	}

	var jsonArray map[string]interface{}
	err = json.Unmarshal(body, &jsonArray)
	if err != nil {
		fmt.Println(err)

	}

	var tabs []Tab
	for _, item := range jsonArray {
		if innerMap, ok := item.(map[string]interface{}); ok {
			// check if each field exists and is not nil
			id, okID := innerMap["id"].(string)
			downloadLink, okDownloadLink := innerMap["downloadLink"].(string)
			title, okTitle := innerMap["title"].(string)
			artist, okArtist := innerMap["artist"].(string)

			if !okID || !okDownloadLink || !okTitle || !okArtist {
				continue
			}

			tab := Tab{
				id:           id,
				downloadLink: downloadLink,
				title:        title,
				artist:       artist,
			}

			tabs = append(tabs, tab)
		}
	}

	return tabs
}

func downloadTab(url string, artist string, title string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	splitURL := strings.Split(url, ".")
	extension := splitURL[len(splitURL)-1]
	filename := fmt.Sprintf("%s - %s.%s", artist, title, extension)
	out, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	startServerInBackground()

	myApp := app.New()
	myWindow := myApp.NewWindow("Tab Search")
	appName := widget.NewLabel("TabStop")
	appName.Alignment = fyne.TextAlign(fyne.TextAlignCenter)
	appName.TextStyle = fyne.TextStyle{Bold: true}

	input := widget.NewEntry()
	input.SetPlaceHolder("Search...")

	var tabs []Tab

	list := widget.NewList(
		func() int {
			return len(tabs)
		},
		func() fyne.CanvasObject {

			return container.NewBorder(nil, nil,
				widget.NewLabel("template"),
				widget.NewButton("Download", func() {}),
			)
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			tab := tabs[i]
			o.(*fyne.Container).Objects[0].(*widget.Label).SetText(fmt.Sprintf("%s - %s", tab.title, tab.artist))
			o.(*fyne.Container).Objects[1].(*widget.Button).OnTapped = func(t Tab) func() {
				return func() {
					err := downloadTab(t.downloadLink, t.artist, t.title)
					if err != nil {
						fmt.Println("Download failed:", err)
					} else {
						fmt.Println("Downloaded:", t.title)
					}
				}
			}(tab)
		})

	settingsContent := container.NewBorder(nil, nil, nil, nil)
	myCanvas := myWindow.Canvas()
	myCanvas.SetContent(settingsContent)
	thePopup := widget.NewModalPopUp(myCanvas.Content(), myCanvas)

	settingsIcon := widget.NewButtonWithIcon("", theme.SettingsIcon(), func() {
		log.Println("tapped settings")
		widget.ShowModalPopUp(thePopup, myCanvas) // show modal popup
	})
	topBar := container.NewBorder(nil, nil, appName, settingsIcon)

	searchContainer := container.NewBorder(topBar, nil, nil, widget.NewButton("Search", func() {
		tabs = getTabs(input.Text)
		list.Refresh()
	}), input)

	content := container.NewBorder(searchContainer, nil, nil, nil, list)
	myWindow.SetContent(content)
	myWindow.Resize(fyne.NewSize(400, 600))
	myWindow.ShowAndRun()
}
