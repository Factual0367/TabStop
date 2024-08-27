package ui

import (
	"TabStop/utils"
	"fmt"

	"fyne.io/fyne/theme"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func Run() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Tab Search")
	appName := widget.NewLabel("TabStop")
	appName.Alignment = fyne.TextAlignCenter
	appName.TextStyle = fyne.TextStyle{Bold: true}

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
						o.(*fyne.Container).Objects[1].(*widget.Button).Icon = theme.ConfirmIcon()
					}
				}
			}(tab)
		})

	settingsIcon := widget.NewButtonWithIcon("", theme.SettingsIcon(),
		func() { utils.ShowSettings(myWindow) })

	topBar := container.NewBorder(nil, nil, appName, settingsIcon)

	searchContainer := container.NewBorder(topBar, nil, nil, widget.NewButton("Search", func() {
		tabs = utils.GetTabs(input.Text)
		list.Refresh()
	}), input)

	content := container.NewBorder(searchContainer, nil, nil, nil, list)
	myWindow.SetContent(content)
	myWindow.Resize(fyne.NewSize(400, 600))
	myWindow.ShowAndRun()
}
