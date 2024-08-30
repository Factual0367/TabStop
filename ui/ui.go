package ui

import (
	"TabStop/utils"
	"fmt"
	"path"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/skratchdot/open-golang/open"
)

func createSearchTab() *container.TabItem {
	input := widget.NewEntry()
	input.SetPlaceHolder("Search...")

	var tabs []utils.Tab

	list := widget.NewList(
		func() int { return len(tabs) },
		func() fyne.CanvasObject {
			return container.NewBorder(nil, nil,
				widget.NewLabel("template"),
				widget.NewButtonWithIcon("", theme.DownloadIcon(), func() {}),
			)
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			tab := tabs[i]
			o.(*fyne.Container).Objects[0].(*widget.Label).SetText(fmt.Sprintf("%s - %s", tab.Title, tab.Artist))
			o.(*fyne.Container).Objects[1].(*widget.Button).OnTapped = func(t utils.Tab) func() {
				return func() {
					err := utils.DownloadTab(t.DownloadLink, t.Artist, t.Title)
					if err != nil {
						fmt.Println("Download failed:", err)
						o.(*fyne.Container).Objects[1].(*widget.Button).Icon = theme.ErrorIcon()
					} else {
						fmt.Println("Downloaded:", t.Title)
						o.(*fyne.Container).Objects[1].(*widget.Button).Icon = theme.FolderOpenIcon()
						o.(*fyne.Container).Objects[1].(*widget.Button).OnTapped = func() {
							downloadDir := utils.GetCurrentDownloadFolder()
							err := open.Run(downloadDir)
							if err != nil {
								return
							}
						}
					}
				}
			}(tab)
		})

	searchBtn := widget.NewButton("Search", func() {
		tabs = utils.GetTabs(input.Text)
		list.Refresh()
	})

	input.OnSubmitted = func(text string) {
		tabs = utils.GetTabs(text)
		list.Refresh()
	}

	rightSideButtons := container.NewHBox(searchBtn)
	searchContainer := container.NewBorder(nil, nil, nil, rightSideButtons, input)

	searchContent := container.NewBorder(searchContainer, nil, nil, nil, list)

	return container.NewTabItem("Search", searchContent)
}

func createMyTabsTab() *container.TabItem {
	downloadDir := utils.GetCurrentDownloadFolder()
	savedTabs := utils.GetSavedTabs()

	list := widget.NewList(
		func() int { return len(savedTabs) },
		func() fyne.CanvasObject {
			return container.NewBorder(nil, nil,
				widget.NewLabel("template"),
				widget.NewButtonWithIcon("", theme.FileIcon(), func() {}),
			)
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			filename := savedTabs[i]
			filepath := path.Join(downloadDir, filename)
			o.(*fyne.Container).Objects[0].(*widget.Label).SetText(fmt.Sprintf("%s", filename))
			o.(*fyne.Container).Objects[1].(*widget.Button).OnTapped = func() {
				err := open.Run(filepath)
				if err != nil {
					fmt.Println("Failed to open file:", err)
					o.(*fyne.Container).Objects[1].(*widget.Button).Icon = theme.ErrorIcon()
				}
			}
		})

	content := container.NewBorder(nil, nil, nil, nil, list)
	return container.NewTabItem("My Tabs", content)
}

func createSettingsTab() *container.TabItem {
	downloadDir := utils.GetCurrentDownloadFolder()
	currentDownloadDirMsg := fmt.Sprintf("Current download directory: \n%s", downloadDir)
	downloadDirLabel := widget.NewLabel(currentDownloadDirMsg)

	changeDownloadLocation := widget.NewButtonWithIcon("Change Download Location", theme.FolderIcon(), func() {
		newDownloadDir := utils.GetFolder()
		downloadDirLabel.SetText(fmt.Sprintf("Current download directory: \n%s", newDownloadDir))
	})

	padding := widget.NewLabel("")

	settingsInnerContent := container.NewGridWithRows(4, padding, downloadDirLabel, changeDownloadLocation)
	settingsContent := container.NewVBox(
		padding,
		container.NewGridWithColumns(3, padding, settingsInnerContent, padding),
	)
	settingsContentBordered := container.NewBorder(nil, nil, nil, nil, settingsContent)

	return container.NewTabItem("Settings", settingsContentBordered)
}

func Run() {
	myApp := app.New()
	myApp.Settings().SetTheme(theme.DarkTheme())
	myWindow := myApp.NewWindow("Tab Search")

	searchTab := createSearchTab()
	searchTab.Icon = theme.SearchIcon()
	myTabsTab := createMyTabsTab()
	myTabsTab.Icon = theme.StorageIcon()
	settingsTab := createSettingsTab()
	settingsTab.Icon = theme.SettingsIcon()
	appTabs := container.NewAppTabs(
		searchTab,
		myTabsTab,
		settingsTab,
	)

	// required to refresh the myTabs list every time
	// the tab gets selected
	appTabs.OnSelected = func(tab *container.TabItem) {
		if tab.Text == "My Tabs" {
			newMyTabsTab := createMyTabsTab()
			newMyTabsTab.Icon = theme.StorageIcon()
			appTabs.Items[1] = newMyTabsTab
			appTabs.Refresh()
		}
	}

	appTabs.SetTabLocation(container.TabLocationLeading)
	myWindow.SetContent(appTabs)
	myWindow.Resize(fyne.NewSize(400, 600))
	myWindow.ShowAndRun()
}
