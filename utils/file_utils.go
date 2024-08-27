package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/sqweek/dialog"
)

type Tab struct {
	ID           string
	DownloadLink string
	Title        string
	Artist       string
}

func Exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func DownloadTab(url string, artist string, title string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	splitURL := strings.Split(url, ".")
	extension := splitURL[len(splitURL)-1]

	filename := fmt.Sprintf("%s - %s.%s", artist, title, extension)
	homedir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println(err)
	}
	fileLocation := path.Join(homedir, ".tabStop")
	tabStopCfgExists, err := Exists(fileLocation)
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

func SaveDownloadLocation(downloadDir string) {
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

func GetTabs(query string) []Tab {
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

	var tabs []Tab
	for _, item := range jsonArray {
		if innerMap, ok := item.(map[string]interface{}); ok {
			id, okID := innerMap["id"].(string)
			downloadLink, okDownloadLink := innerMap["downloadLink"].(string)
			title, okTitle := innerMap["title"].(string)
			artist, okArtist := innerMap["artist"].(string)

			if !okID || !okDownloadLink || !okTitle || !okArtist {
				continue
			}

			tab := Tab{
				ID:           id,
				DownloadLink: downloadLink,
				Title:        title,
				Artist:       artist,
			}

			tabs = append(tabs, tab)
		}
	}
	return tabs
}

func ShowSettings(w fyne.Window) (modal *widget.PopUp) {
	modal = widget.NewModalPopUp(
		container.NewVBox(
			widget.NewLabel("Set Download Location"),
			widget.NewButtonWithIcon("", theme.FolderIcon(), func() { GetFolder() }),
			widget.NewButton("Close", func() { modal.Hide() }),
		),
		w.Canvas(),
	)
	modal.Show()
	return modal
}

func GetFolder() {
	directory, err := dialog.Directory().Title("Select Folder").Browse()
	if err != nil {
		fmt.Println(err)
	}
	SaveDownloadLocation(directory)
}
