package main

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

var data = [][]string{
	{"Song ID", "Song"},
	{"bottom left", "bottom right"},
}

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("TabStop")

	input := widget.NewEntry()
	input.SetPlaceHolder("Enter text...")

	list := widget.NewTable(
		func() (int, int) {
			return len(data), len(data[0])
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("wide content")
		},
		func(i widget.TableCellID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(data[i.Row][i.Col])
		})

	// Function to update the data and refresh the table
	updateData := func(query string) {
		// Example logic to populate the table with new data
		// based on the search query
		data = [][]string{
			{"Song ID", "Song"},
			{"1", "Song One"},
			{"2", "Song Two"},
			{"3", "Song Three"},
		}

		// Refresh the table to show the updated data
		list.Refresh()
	}

	content := container.NewVBox(
		input,
		widget.NewButton("Search", func() {
			log.Println("Search query:", input.Text)
			updateData(input.Text) // Call updateData with the search query
		}),
		list,
	)

	myWindow.SetContent(content)
	myWindow.ShowAndRun()
}
