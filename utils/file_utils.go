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

func GetCurrentDownloadFolder() string {
	var downloadDir string
	homedir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println(err)
		return ""
	}
	fileLocation := path.Join(homedir, ".tabStop")

	tabStopCfgExists, err := Exists(fileLocation)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	if tabStopCfgExists {
		file, err := os.ReadFile(fileLocation)
		if err != nil {
			fmt.Println(err)
			return ""
		}

		userData := map[string]interface{}{}
		err = json.Unmarshal(file, &userData)
		if err != nil {
			fmt.Println(err)
			return ""
		}

		if dl, exists := userData["downloadLocation"].(string); exists {
			downloadDir = dl
		} else {
			downloadDir, err = os.Getwd()
			if err != nil {
				fmt.Println(err)
				return ""
			}
		}

	} else {
		downloadDir, err = os.Getwd()
		if err != nil {
			fmt.Println(err)
			return ""
		}
	}

	return downloadDir
}

func GetSavedTabs() (map[string]string, error) {
	homedir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get user home directory: %w", err)
	}

	fileLocation := path.Join(homedir, ".tabStop")

	userCfgExists, err := Exists(fileLocation)
	if err != nil {
		return nil, fmt.Errorf("failed to check if config file exists: %w", err)
	}

	if !userCfgExists {
		return nil, fmt.Errorf("no configuration file found")
	}

	file, err := os.ReadFile(fileLocation)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	userData := map[string]interface{}{}
	err = json.Unmarshal(file, &userData)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal config file: %w", err)
	}

	downloadedTabs, ok := userData["downloadedTabs"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("no downloaded tabs found")
	}

	tabs := make(map[string]string)
	for key, value := range downloadedTabs {
		if strValue, ok := value.(string); ok {
			tabs[key] = strValue
		}
	}

	return tabs, nil
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

	downloadDir := GetCurrentDownloadFolder()

	filepath := path.Join(downloadDir, filename)

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	SaveDownloadLocation(downloadDir, filename)

	return nil
}

func SaveDownloadLocation(downloadDir string, filename string) {
	homedir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println(err)
	}
	fileLocation := path.Join(homedir, ".tabStop")

	userData := map[string]interface{}{}

	// check if file exists
	userCfgExists, err := Exists(fileLocation)

	if userCfgExists {
		file, _ := os.ReadFile(fileLocation)
		json.Unmarshal(file, &userData)

		if _, exists := userData["downloadedTabs"]; !exists {
			userData["downloadedTabs"] = map[string]interface{}{}
		}

		downloadedTabs := userData["downloadedTabs"].(map[string]interface{})
		downloadedTabs[filename] = path.Join(downloadDir, filename)

	} else {
		userData = map[string]interface{}{
			"downloadLocation": downloadDir,
			"downloadedTabs": map[string]interface{}{
				filename: path.Join(downloadDir, filename),
			},
		}
	}

	fileData, _ := json.MarshalIndent(userData, "", "  ")
	_ = os.WriteFile(fileLocation, fileData, 0644)
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

	downloadDir := GetCurrentDownloadFolder()
	currentDownloadDirMsg := fmt.Sprintf("Current download directory: \n%s", downloadDir)
	modal = widget.NewModalPopUp(
		container.NewVBox(
			widget.NewLabel(currentDownloadDirMsg),
			widget.NewButtonWithIcon("Change", theme.FolderIcon(), func() { GetFolder() }),
			widget.NewButtonWithIcon("Close", theme.WindowCloseIcon(), func() { modal.Hide() }),
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
	SaveDownloadLocation(directory, "")
}
