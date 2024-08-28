package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strings"

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

func GetSavedTabs() []string {
	downloadDir := GetCurrentDownloadFolder()

	savedTabs, _ := os.ReadDir(downloadDir)

	tabs := make([]string, 0, len(savedTabs))
	for _, e := range savedTabs {
		tabs = append(tabs, e.Name())
	}

	return tabs

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

	return nil
}

func SaveDownloadLocation(downloadDir string) {
	homedir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println(err)
	}
	fileLocation := path.Join(homedir, ".tabStop")

	userData := map[string]interface{}{}

	userCfgExists, err := Exists(fileLocation)

	if userCfgExists {
		file, _ := os.ReadFile(fileLocation)
		json.Unmarshal(file, &userData)

		userData["downloadLocation"] = downloadDir

	} else {
		userData = map[string]interface{}{
			"downloadLocation": downloadDir,
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

func GetFolder() string {
	directory, err := dialog.Directory().Title("Select Folder").Browse()
	if err != nil {
		fmt.Println(err)
	}
	SaveDownloadLocation(directory)
	return directory
}
