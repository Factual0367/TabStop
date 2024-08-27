package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strings"

	"fyne.io/fyne/theme"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/gin-gonic/gin"
	"github.com/onurhanak/songsterrapi"
	"github.com/sqweek/dialog"
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

	query = strings.ReplaceAll(query, " ", "%20")
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
	fmt.Println(jsonArray)
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
	fmt.Println(tabs)
	return tabs
}

// exists returns whether the given file or directory exists
func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
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
	// check if custom download location exists
	// if exists save tabs there
	homedir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println(err)
	}
	fileLocation := path.Join(homedir, ".tabStop")
	tabStopCfgExists, err := exists(fileLocation)
	if err != nil {
		fmt.Println(err)
	}
	if tabStopCfgExists {
		downloadDir, err := os.ReadFile(fileLocation)
		if err != nil {
			fmt.Println(err)
		}
		filename = path.Join(string(downloadDir), filename)

	}

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

func saveDownloadLocation(downloadDir string) {
	homedir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println(err)
	}
	fileLocation := path.Join(homedir, ".tabStop")
	f, err := os.Create(fileLocation)
	if err != nil {
		fmt.Println(err)
	}
	l, err := f.WriteString(downloadDir)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(l, "bytes written.")
	err = f.Close()
	if err != nil {
		fmt.Println(err)
	}
}
func getFolder() {
	directory, err := dialog.Directory().Title("Select Folder").Browse()
	if err != nil {

		fmt.Println(err)

	}

	saveDownloadLocation(directory)
}

func showSettings(w fyne.Window) (modal *widget.PopUp) {

	modal = widget.NewModalPopUp(
		container.NewVBox(
			widget.NewLabel("Set Download Location"),
			widget.NewButtonWithIcon("", theme.FolderIcon(), func() { getFolder() }),
			widget.NewButton("Close", func() { modal.Hide() }),
		),
		w.Canvas(),
	)

	modal.Show()
	return modal
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
				widget.NewButtonWithIcon("", theme.DownloadIcon(), func() {}),
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
						o.(*fyne.Container).Objects[1].(*widget.Button).Icon = theme.ErrorIcon()

					} else {
						fmt.Println("Downloaded:", t.title)
						o.(*fyne.Container).Objects[1].(*widget.Button).Icon = theme.ConfirmIcon()
					}
				}
			}(tab)
		})

	settingsIcon := widget.NewButtonWithIcon("", theme.SettingsIcon(),

		func() { showSettings(myWindow) })

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
